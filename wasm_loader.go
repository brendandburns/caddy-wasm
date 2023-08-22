package caddy_wasm

import (
	"context"
	"sync"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

type ModuleInstance struct {
	module api.Module
	mut    sync.Mutex
}

func (m *ModuleInstance) acquire() {
	m.mut.Lock()
}

func (m *ModuleInstance) tryAcquire() bool {
	return m.mut.TryLock()
}

func (m *ModuleInstance) Release() {
	m.mut.Unlock()
}

func (m *ModuleInstance) Module() api.Module {
	return m.module
}

type WebAssemblyLoader struct {
	mut       sync.RWMutex
	instances []*ModuleInstance

	rt           wazero.Runtime
	module       wazero.CompiledModule
	moduleConfig wazero.ModuleConfig
}

func (w *WebAssemblyLoader) Get() *ModuleInstance {
	w.mut.RLock()
	defer w.mut.RUnlock()

	for ix := range w.instances {
		if w.instances[ix].tryAcquire() {
			return w.instances[ix]
		}
	}
	return nil
}

func (w *WebAssemblyLoader) GetOrLoad(ctx context.Context) (*ModuleInstance, error) {
	// try to use the ones we have
	if inst := w.Get(); inst != nil {
		return inst, nil
	}

	// ok we need to make a new one
	mod, err := w.rt.InstantiateModule(ctx, w.module, w.moduleConfig)
	if err != nil {
		return nil, err
	}
	instance := &ModuleInstance{
		module: mod,
		mut:    sync.Mutex{},
	}
	instance.acquire()

	w.mut.Lock()
	defer w.mut.Unlock()
	w.instances = append(w.instances, instance)
	return instance, nil
}
