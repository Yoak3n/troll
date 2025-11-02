declare global {
  interface Window {
    $loadingBar?: ReturnType<typeof import('naive-ui')['useLoadingBar']>
    $dialog?: ReturnType<typeof import('naive-ui')['useDialog']>
    $message?: ReturnType<typeof import('naive-ui')['useMessage']>
    $notification?: ReturnType<typeof import('naive-ui')['useNotification']>
  }
}

export {}
