package disgopage

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	manager *paginatorManager = newManager()
)

// paginatorManager is the main controller for the paginator. It contains the
// configuration used by all paginators, as well as all active
// paginators
type paginatorManager struct {
	mutex      sync.Mutex
	paginators map[string]*Paginator
}

// newManager creates a new paginator manager. It takes a variadic number of ConfigOpt
// options to configure the paginator. The default configuration is used if no
// options are provided.  The manager starts a goroutine to clean up expired paginators.
// The cleanup goroutine runs every minute and removes expired paginators.
func newManager() *paginatorManager {
	manager := &paginatorManager{
		paginators: map[string]*Paginator{},
		mutex:      sync.Mutex{},
	}
	manager.startCleanup()

	return manager
}

// addPaginator adds a paginator to the manager.
func (m *paginatorManager) addPaginator(paginator *Paginator) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.paginators[paginator.id] = paginator
	log.WithFields(log.Fields{"paginator": paginator.id, "count": len(m.paginators)}).Debug("added paginator to manager")
}

// removePaginator removes a paginator from the manager and performs any necessary cleanup.
// It contains the shared logic used by `Remove“ and `cleanup`.
func (m *paginatorManager) removePaginator(p *Paginator) {
	delete(m.paginators, p.id)
	log.WithFields(log.Fields{"paginator": p.id}).Debug("removed paginator from manager")
}

// cleanup removes expired paginators from the manager.
func (m *paginatorManager) cleanup() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, p := range m.paginators {
		p.cleanup()
	}
}

// startCleanup starts a goroutine that cleans up expired paginators every minute.
func (m *paginatorManager) startCleanup() {
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			log.WithFields(log.Fields{"count": len(m.paginators)}).Trace("running paginator cleanup")
			m.cleanup()
		}
	}()
}
