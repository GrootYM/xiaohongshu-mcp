package browser

import (
	"sync"

	"github.com/go-rod/rod"
	"github.com/sirupsen/logrus"
	"github.com/xpzouying/headless_browser"
)

// Pool 浏览器池，复用浏览器实例以提升性能
type Pool struct {
	browser  *headless_browser.Browser
	headless bool
	binPath  string
	mu       sync.Mutex
}

var (
	globalPool *Pool
	poolOnce   sync.Once
)

// InitPool 初始化全局浏览器池（服务启动时调用）
func InitPool(headless bool, binPath string) {
	poolOnce.Do(func() {
		globalPool = &Pool{
			headless: headless,
			binPath:  binPath,
		}
		// 预热浏览器
		if err := globalPool.warmup(); err != nil {
			logrus.Warnf("浏览器预热失败: %v", err)
		}
	})
}

// GetPool 获取全局浏览器池
func GetPool() *Pool {
	return globalPool
}

// warmup 预热浏览器实例
func (p *Pool) warmup() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.browser != nil {
		return nil
	}

	logrus.Info("正在预热浏览器...")
	p.browser = p.createBrowser()
	logrus.Info("浏览器预热完成")
	return nil
}

// createBrowser 创建新的浏览器实例
func (p *Pool) createBrowser() *headless_browser.Browser {
	return NewBrowser(p.headless, WithBinPath(p.binPath))
}

// GetPage 获取一个可用的页面
func (p *Pool) GetPage() (*rod.Page, func(), error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// 如果浏览器不存在或已关闭，重新创建
	if p.browser == nil {
		logrus.Info("创建新的浏览器实例...")
		p.browser = p.createBrowser()
	}

	// 创建新页面
	page := p.browser.NewPage()

	// 返回页面和清理函数（只关闭页面，不关闭浏览器）
	cleanup := func() {
		if page != nil {
			_ = page.Close()
		}
	}

	return page, cleanup, nil
}

// Close 关闭浏览器池
func (p *Pool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.browser != nil {
		p.browser.Close()
		p.browser = nil
	}
}

// Restart 重启浏览器（当浏览器出问题时调用）
func (p *Pool) Restart() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.browser != nil {
		p.browser.Close()
	}
	logrus.Info("重启浏览器...")
	p.browser = p.createBrowser()
}
