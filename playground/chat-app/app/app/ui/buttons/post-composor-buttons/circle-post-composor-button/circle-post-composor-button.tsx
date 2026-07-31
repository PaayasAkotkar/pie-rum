'use client'
import { motion, } from "framer-motion"
import { useState, useEffect, } from "react"
import { $defaultBlurCap, StyleSheetNode, motionDefaultProps } from "../../../../misc/types"
import { useResponsive } from "../../../../services/cssx/responsive/services/responsive/use-responsive"
import BlurFilter from "../../../../services/filter/blur/blur"
import NormalFilter from "../../../../services/filter/normal/normal"

interface input {
    src?: string,
    title?: string,
    alt?: string
    type?: "button" | "submit" | "reset"
    node?: motionDefaultProps<HTMLButtonElement>
    pause?: boolean
    trigger?: boolean
    styleSheetNode?: StyleSheetNode
    delay?: number
    children?: React.ReactNode
}
export default function CirclePostComposorButton({
    src, title, alt, node, pause, type = "button",
    styleSheetNode, delay, trigger,
    children
}: input) {

    const [click, setClick] = useState(pause ?? false)

    useEffect(() => {
        if (pause !== undefined) setClick(pause)
    }, [pause])

    const handleCombinedClick = (e: React.MouseEvent<HTMLButtonElement>) => {
        e.preventDefault();
        e.stopPropagation();

        if (!click) setClick(true);

        const currentTarget = e.currentTarget;
        const rect = currentTarget.getBoundingClientRect();

        if (node?.onClick) {
            setTimeout(() => {
                const mockEvent = {
                    type: 'click',
                    target: currentTarget,
                    currentTarget: {
                        ...currentTarget,
                        getBoundingClientRect: () => rect,
                    },
                    clientX: e.clientX,
                    clientY: e.clientY,
                    defaultPrevented: true,
                    stopPropagation: () => { },
                    preventDefault: () => { },
                } as unknown as React.MouseEvent<HTMLButtonElement>;

                node.onClick?.(mockEvent);
            }, 200);
        }
    };
    const { font, circle } = useResponsive()
    const _size = styleSheetNode?.size ?? 20
    const f = styleSheetNode?.btnFontSize ?? 20
    const tc = styleSheetNode?.btnTextColor ?? 'var(--background)'
    const be = styleSheetNode?.blurEffect ?? false
    const bc = styleSheetNode?.blurCap ?? $defaultBlurCap
    const w = circle(_size)
    const h = circle(_size)
    const fontSize = font(f)

    useEffect(() => {
        if (!trigger) return
        setClick(trigger)
    }, [trigger])
    return (

        <motion.button
            type={type ?? "button"}
            {...node}
            onClick={handleCombinedClick}
            initial="idle"
            style={{
                width: w,
                height: h,
                outline: 'none',
                border: 'none',
                backgroundColor: styleSheetNode?.btnTheme ?? "var(--foreground)",
                backgroundImage: styleSheetNode?.imgTheme ?? "none",
                padding: 0,
                flexShrink: 0,
            }}
            className="relative rounded-full flex items-center justify-center overflow-hidden border-none cursor-pointer"
        >
            {click && !be && (
                <NormalFilter direction={styleSheetNode?.filterDirection} delay={delay} onComplete={() => setClick(false)} tricolor={styleSheetNode?.triColor} triggerAnimation={click} />
            )}
            {click && be && (
                <BlurFilter blurCap={bc} delay={delay} onComplete={() => setClick(false)} tricolor={styleSheetNode?.triColor} triggerAnimation={click} />
            )}

            <div className="absolute inset-0 flex items-center justify-center z-20 pointer-events-none p-[15%]">
                {src && (
                    <img
                        className="w-full h-full object-contain"
                        src={src}
                        alt={alt ?? 'button icon'}
                    />
                )}
                {title && (
                    <span
                        className="whitespace-nowrap"
                        style={{
                            color: tc,
                            fontSize: fontSize
                        }}
                    >
                        {title}
                    </span>
                )}
                {children && children}
            </div>
        </motion.button>

    )
}