export const MOBILE_WORKBENCH_NAVIGATE_EVENT = 'nova:mobile-workbench:navigate'

export function requestMobileWorkbenchDestination(destinationId: string) {
  if (typeof window === 'undefined') return
  window.dispatchEvent(new CustomEvent(MOBILE_WORKBENCH_NAVIGATE_EVENT, { detail: { destinationId } }))
}

export function mobileWorkbenchDestinationFromEvent(event: Event): string | null {
  if (!(event instanceof CustomEvent)) return null
  const detail = event.detail as { destinationId?: unknown } | null
  return typeof detail?.destinationId === 'string' ? detail.destinationId : null
}
