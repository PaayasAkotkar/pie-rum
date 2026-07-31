import { useState, useCallback } from 'react'

interface input {
    timeout?: number
    onSuccess?: () => void
    onError?: (error: Error) => void
}

interface output {
    isCopied: boolean
    copy: (text: string) => Promise<void>
    reset: () => void
    error: Error | null
}
// useCopy is a custom hook that provides copy to clipboard functionality
export function useCopy(options: input = {}): output {
    const { timeout = 1000, onSuccess, onError } = options

    const [isCopied, setIsCopied] = useState(false)
    const [error, setError] = useState<Error | null>(null)

    const copy = useCallback(
        async (text: string) => {
            setError(null)

            try {
                if (navigator.clipboard && window.isSecureContext) {
                    await navigator.clipboard.writeText(text)
                }

                setIsCopied(true)
                onSuccess?.()

                setTimeout(() => {
                    setIsCopied(false)
                }, timeout)

            } catch (err) {
                const copyError = err instanceof Error ? err : new Error('Failed to copy')
                setError(copyError)
                onError?.(copyError)
                setIsCopied(false)
            }
        },
        [timeout, onSuccess, onError]
    )

    const reset = useCallback(() => {
        setIsCopied(false)
        setError(null)
    }, [])

    return {
        isCopied,
        copy,
        reset,
        error,
    }
}