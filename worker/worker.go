package worker

import (
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Job struct {
	MessageID string // used as the correlation key for status inquiry
	FileName  string // e.g. PACS0081720000000001.xml
	Content   []byte // the MX (ISO 20022) XML payload
	MsgType   string // sub-folder under the output dir, e.g. "pacs008"
}

const (
	StatusPending    = "PENDING"
	StatusProcessing = "PROCESSING"
	StatusCompleted  = "COMPLETED"
	StatusFailed     = "FAILED"
)

type JobResult struct {
	MessageID   string    `json:"message_id"`
	MsgType     string    `json:"message_type"`
	Status      string    `json:"status"`
	FilePath    string    `json:"file_path,omitempty"`
	Error       string    `json:"error,omitempty"`
	SubmittedAt time.Time `json:"submitted_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Pool struct {
	jobs      chan Job
	results   map[string]*JobResult
	mu        sync.RWMutex
	outputDir string
	wg        sync.WaitGroup
}

// NewPool creates the output directory (if needed) and starts
// workerCount goroutines listening on a channel of size queueSize.
func NewPool(workerCount int, outputDir string, queueSize int) *Pool {
	p := &Pool{
		jobs:      make(chan Job, queueSize),
		results:   make(map[string]*JobResult),
		outputDir: outputDir,
	}
	_ = os.MkdirAll(outputDir, 0755)

	for i := 0; i < workerCount; i++ {
		p.wg.Add(1)
		go p.startWorker(i)
	}
	return p
}

func (p *Pool) startWorker(_ int) {
	defer p.wg.Done()
	for job := range p.jobs {
		p.setStatus(job.MessageID, job.MsgType, StatusProcessing, "", "")

		typeDir := filepath.Join(p.outputDir, job.MsgType)
		if err := os.MkdirAll(typeDir, 0755); err != nil {
			p.setStatus(job.MessageID, job.MsgType, StatusFailed, "", err.Error())
			continue
		}

		fullPath := filepath.Join(typeDir, job.FileName)
		if err := os.WriteFile(fullPath, job.Content, 0644); err != nil {
			p.setStatus(job.MessageID, job.MsgType, StatusFailed, "", err.Error())
			continue
		}

		p.setStatus(job.MessageID, job.MsgType, StatusCompleted, fullPath, "")
	}
}

func (p *Pool) setStatus(id, msgType, status, path, errMsg string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	existing, ok := p.results[id]
	now := time.Now()
	if !ok {
		existing = &JobResult{MessageID: id, MsgType: msgType, SubmittedAt: now}
	}
	existing.Status = status
	existing.MsgType = msgType
	existing.UpdatedAt = now
	if path != "" {
		existing.FilePath = path
	}
	existing.Error = errMsg
	p.results[id] = existing
}

func (p *Pool) Submit(job Job) {
	p.mu.Lock()
	p.results[job.MessageID] = &JobResult{
		MessageID:   job.MessageID,
		MsgType:     job.MsgType,
		Status:      StatusPending,
		SubmittedAt: time.Now(),
		UpdatedAt:   time.Now(),
	}
	p.mu.Unlock()

	p.jobs <- job
}

func (p *Pool) GetStatus(id string) (*JobResult, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	r, ok := p.results[id]
	if !ok {
		return nil, false
	}
	cp := *r
	return &cp, true
}

func (p *Pool) Shutdown() {
	close(p.jobs)
	p.wg.Wait()
}
