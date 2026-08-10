import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { QueryClientProvider } from '@tanstack/react-query'
import { ThemeProvider } from 'next-themes'
import { afterEach, beforeEach, describe, expect, it } from 'vitest'
import { http, HttpResponse } from 'msw'
import { server } from '@/test/msw/server'
import { queryClient } from '@/lib/query-client'
import { TooltipProvider } from '@/components/ui/tooltip'
import { Toaster } from '@/components/ui/sonner'
import { setConfiguredLocale } from '@/i18n'
import App from '@/App'

const WORKSPACE_ONE = '/books/one'
const WORKSPACE_TWO = '/books/two'

function renderApp() {
  return render(
    <QueryClientProvider client={queryClient}>
      <ThemeProvider attribute="data-theme" defaultTheme="dark" enableSystem themes={['light', 'dark']}>
        <TooltipProvider>
          <App />
          <Toaster richColors closeButton />
        </TooltipProvider>
      </ThemeProvider>
    </QueryClientProvider>,
  )
}

function workspaceHandlers(saveRequests: unknown[], switchRequests: unknown[]) {
  return [
    http.get('/api/workspace/current', () =>
      HttpResponse.json({ workspace: WORKSPACE_ONE, has_state: true }),
    ),
    http.get('/api/books', () =>
      HttpResponse.json({
        sort_mode: 'recent',
        books: [
          { name: '创作项目一', path: WORKSPACE_ONE, author: '', last_opened_at: '' },
          { name: '创作项目二', path: WORKSPACE_TWO, author: '', last_opened_at: '' },
        ],
      }),
    ),
    http.get('/api/workspace/tree', () => HttpResponse.json([])),
    http.get('/api/workspace/summary', () =>
      HttpResponse.json({
        title: '创作项目一',
        author: '',
        chapter_count: 1,
        total_words: 4,
        chapters: [
          {
            path: 'chapters/ch01.md',
            file_name: 'ch01.md',
            display_title: '第一章',
            index: 1,
            words: 4,
            status: 'draft',
            confirmed: false,
            updated_at: '2026-08-11T00:00:00Z',
            volume: '',
            volume_path: '',
          },
        ],
        chapter_plans: [],
      }),
    ),
    http.get('/api/workspace/file', ({ request }) => {
      const path = new URL(request.url).searchParams.get('path') || 'chapters/ch01.md'
      return HttpResponse.json({
        workspace: WORKSPACE_ONE,
        path,
        content: path === 'chapters/ch01.md' ? '第一章' : '',
        revision: 'r1',
      })
    }),
    http.post('/api/workspace/file', async ({ request }) => {
      const body = await request.json()
      saveRequests.push(body)
      return HttpResponse.json({ error: '磁盘写入失败' }, { status: 500 })
    }),
    http.post('/api/workspace/switch', async ({ request }) => {
      const body = await request.json()
      switchRequests.push(body)
      return HttpResponse.json({ workspace: WORKSPACE_TWO, message: 'ok' })
    }),
    http.get('/api/tasks', () =>
      HttpResponse.json({ action_required_count: 0, tasks: [] }),
    ),
    http.get('/api/image-presets', () => HttpResponse.json({ presets: [] })),
    http.get('/api/workspace/document-review', () =>
      HttpResponse.json({ workspace: WORKSPACE_ONE, review_thread: { id: '', comments: [] } }),
    ),
    http.get('/api/update/check', () =>
      HttpResponse.json({
        current_version: '0.3.3',
        latest_version: '0.3.3',
        update_available: false,
        can_install: false,
        platform: 'windows',
        release_url: '',
        published_at: '',
      }),
    ),
  ]
}

describe('应用级项目切换保存保护', () => {
  let restoreMatchMedia: () => void

  beforeEach(() => {
    const originalMatchMedia = window.matchMedia
    Object.defineProperty(window, 'matchMedia', {
      configurable: true,
      writable: true,
      value: (query: string) => ({
        matches: query.includes('1099px') || query.includes('639px'),
        media: query,
        onchange: null,
        addListener: () => {},
        removeListener: () => {},
        addEventListener: () => {},
        removeEventListener: () => {},
        dispatchEvent: () => false,
      }),
    })
    restoreMatchMedia = () => {
      Object.defineProperty(window, 'matchMedia', {
        configurable: true,
        writable: true,
        value: originalMatchMedia,
      })
    }

    localStorage.clear()
    setConfiguredLocale('zh-CN')
    Object.defineProperty(document, 'visibilityState', { configurable: true, value: 'hidden' })
    localStorage.setItem('nova:content-mode', 'ide')
    localStorage.setItem('nova:onboarding:v1', JSON.stringify({ version: 1, skipped: true }))
    localStorage.setItem(`nova.layout.tabs:${WORKSPACE_ONE}`, JSON.stringify([
      { kind: 'file', path: 'chapters/ch01.md' },
    ]))
    localStorage.setItem(`nova.layout.activeTab:${WORKSPACE_ONE}`, 'file:chapters/ch01.md')
  })

  afterEach(() => {
    restoreMatchMedia()
  })

  it('保存失败时留在当前创作项目并显示明确错误，且不发起项目切换请求', async () => {
    const user = userEvent.setup()
    const saveRequests: unknown[] = []
    const switchRequests: unknown[] = []
    server.use(...workspaceHandlers(saveRequests, switchRequests))

    renderApp()

    const switcher = await screen.findByRole('button', { name: '切换书籍，当前：创作项目一' })
    const editor = await waitFor(() => document.querySelector('.ProseMirror'))
    expect(editor).not.toBeNull()

    await user.click(editor!)
    await user.keyboard('修改')
    await waitFor(() => expect(editor!.textContent).toContain('修改'))

    await user.click(switcher)
    await user.click(await screen.findByRole('menuitem', { name: /创作项目二/ }))

    expect(await screen.findAllByText('保存失败')).not.toHaveLength(0)
    await waitFor(() => expect(saveRequests.length).toBeGreaterThan(0))
    expect(switchRequests).toHaveLength(0)

    await user.keyboard('{Escape}')
    expect(screen.getByRole('button', { name: '切换书籍，当前：创作项目一' })).toBeInTheDocument()
  })

  it('从书架管理切换到其他项目时，保存失败同样留在当前项目并显示错误', async () => {
    const user = userEvent.setup()
    const saveRequests: unknown[] = []
    const switchRequests: unknown[] = []
    server.use(...workspaceHandlers(saveRequests, switchRequests))

    renderApp()

    const switcher = await screen.findByRole('button', { name: '切换书籍，当前：创作项目一' })
    const editor = await waitFor(() => document.querySelector('.ProseMirror'))
    expect(editor).not.toBeNull()

    await user.click(editor!)
    await user.keyboard('修改')
    await waitFor(() => expect(editor!.textContent).toContain('修改'))

    await user.click(switcher)
    await user.click(await screen.findByRole('menuitem', { name: '管理书架' }))
    await user.click(await screen.findByRole('button', { name: /创作项目二/ }))

    expect(await screen.findAllByText('保存失败')).not.toHaveLength(0)
    await waitFor(() => expect(saveRequests.length).toBeGreaterThan(0))
    expect(switchRequests).toHaveLength(0)
    expect(screen.getByRole('button', { name: '切换书籍，当前：创作项目一' })).toBeInTheDocument()
  })
})
