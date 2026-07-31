'use client'
import { AnimatePresence, motion } from "framer-motion"
import { useState } from "react"
import { StyleSheetNode, motionDefaultProps } from "../../../../misc/types"
import { useResponsive } from "../../../../services/cssx/responsive/services/responsive/use-responsive"
import { useBorders } from "@/app/services/cssx/responsive/services/sizes/border"

interface input {
    src?: string
    title?: string
    alt?: string
    node?: motionDefaultProps<HTMLButtonElement>
    styleSheetNode?: StyleSheetNode
    delay?: number
    reverse?: boolean
}

export default function NormalComposorButton({ reverse, delay, src, title, alt, node, styleSheetNode }: input) {
    const [_, setIsHovered] = useState(false);
    const [click, setClick] = useState(false)

    const handleCombinedClick = (e: React.MouseEvent<any>) => {
        if (!click) setClick(true);

        const currentTarget = e.currentTarget;
        const rect = currentTarget.getBoundingClientRect();

        const delayedEvent = {
            ...e,
            currentTarget: { ...currentTarget, getBoundingClientRect: () => rect }
        } as unknown as React.MouseEvent<any>;

        if (node?.onClick) {
            setTimeout(() => node.onClick?.(delayedEvent), delay ?? 400);
        }
    };

    const { clamp, font, square } = useResponsive()
    
    const bw = square(styleSheetNode?.btnWidth ?? 20)
    const bh = square(styleSheetNode?.btnHeight ?? 20)
    
    const { borderWidth } = useBorders()
    const btnRound = clamp(borderWidth, 'px')
    const fontSize = font(styleSheetNode?.btnFontSize ?? 12)

    return (
        <div className="relative flex flex-col items-center justify-end">
            <motion.button
                {...node}
                onClick={handleCombinedClick}
                whileHover={{ scale: 1.04 }}
                whileTap={{ scale: 0.98 }}
                onHoverStart={() => setIsHovered(true)}
                onHoverEnd={() => setIsHovered(false)}
                transition={{ type: "spring", stiffness: 400, damping: 22 }}
                style={{
                    width: bw,
                    height: bh,
                                    outline: 'none',
                border:'none',
                    borderRadius: btnRound,
                    backgroundColor: styleSheetNode?.btnTheme ?? "var(--foreground)",
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: src && title ? 'space-between' : 'center', // Justify between if both exist
                    gap: src && title ? '12px' : '0px',
                    paddingLeft: src && title ? '12px' : '0px',
                    paddingRight: src && title ? '16px' : '0px',
                }}
                className="relative overflow-hidden"
            >
                <AnimatePresence mode="wait">
                    {src && title ? (
                        <div className={`flex w-full items-center justify-between ${reverse ? 'flex-row-reverse' : 'flex-row'}`}>
                            <div className="relative flex-shrink-0" style={{ width: `calc(${bh} * 0.55)`, height: `calc(${bh} * 0.55)` }}>
                                <img
                                    src={src}
                                    alt={alt || "app-icon"}
                                    className="object-contain"
                                />
                            </div>

                            <span
                                className="font-bold text-right whitespace-nowrap leading-none pl-2"
                                style={{
                                    fontSize: fontSize,
                                    color: styleSheetNode?.btnTextColor ?? "black"
                                }}
                            >
                                {title}
                            </span>
                        </div>
                    ) : (
                        !src && title ? (
                            <motion.span
                                key="title-only"
                                initial={{ opacity: 0, y: 4 }}
                                animate={{ opacity: 1, y: 0 }}
                                exit={{ opacity: 0, y: -4 }}
                                className="font-bold px-4 text-center break-words leading-tight"
                                style={{ fontSize: fontSize, color: styleSheetNode?.btnTextColor ?? "black" }}
                            >
                                {title}
                            </motion.span>
                        ) : (
                            src && (
                                <motion.div
                                    key="image-only"
                                    initial={{ opacity: 0, scale: 0.9 }}
                                    animate={{ opacity: 1, scale: 1 }}
                                    exit={{ opacity: 0, scale: 0.9 }}
                                    className="relative w-full h-full flex items-center justify-center"
                                >
                                    <img
                                        src={src}
                                        alt={alt || "app-icon"}
                                        className="object-cover"
                                    />
                                </motion.div>
                            )
                        )
                    )}
                </AnimatePresence>
            </motion.button>
        </div>
    )
}