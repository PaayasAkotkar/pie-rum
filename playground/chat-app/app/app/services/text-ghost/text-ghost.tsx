import { CSSProperties, ReactElement, cloneElement, useRef, useEffect, useCallback } from 'react'

interface input {
    children: ReactElement<any>
    text: string
    fontSize?: string
    lineHeight?: string
    padding?: string
    maxHeight?: string
    containerStyle?: CSSProperties
}

// TextGhost is a custom component which allows the text to be displayed in a ghost which results in having the ms-note alike functionality
// it is ai generated
export function TextGhost({
    children,
    text,
    fontSize = '16px',
    lineHeight = '1.6',
    padding = '4px',
    maxHeight = 'none',
    containerStyle = {},
}: input) {
    const ghostRef = useRef<HTMLDivElement>(null)
    const textareaRef = useRef<HTMLTextAreaElement>(null)

    const mergeRefs = useCallback(
        (node: HTMLTextAreaElement) => {
            textareaRef.current = node

            const { ref } = children as any
            if (typeof ref === 'function') {
                ref(node)
            } else if (ref) {
                ref.current = node
            }
        },
        [children]
    )

    useEffect(() => {
        if (ghostRef.current && textareaRef.current) {
            const ghostHeight = ghostRef.current.scrollHeight
            textareaRef.current.style.height = `${ghostHeight}px`
        }
    }, [text])

    const ghostStyles: CSSProperties = {
        fontSize,
        lineHeight,
        padding,
        whiteSpace: 'pre-wrap',
        wordBreak: 'break-word',
        width: '100%',
        display: 'block',
        visibility: 'hidden',
        position: 'absolute',
        pointerEvents: 'none',
        maxHeight,
        overflow: 'hidden',
        top: 0,
        left: 0,
        zIndex: -1,
        color: 'black',
    }

    const containerStyles: CSSProperties = {
        position: 'relative',
        width: '100%',
        maxHeight,
        ...containerStyle,
    }

    const enhancedChild = cloneElement(children, {
        ref: mergeRefs,
        style: {
            ...(children.props.style || {}),
            height: 'auto',
            minHeight: 'auto',
            boxSizing: 'border-box',
        } as CSSProperties,
    } as any)

    return (
        <div style={containerStyles}>
            <div ref={ghostRef} style={ghostStyles}>
                {(text || ' ') + '\n'}
            </div>
            {enhancedChild}
        </div>
    )
}