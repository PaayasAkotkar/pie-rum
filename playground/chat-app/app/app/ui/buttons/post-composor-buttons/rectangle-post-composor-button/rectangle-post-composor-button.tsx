'use client'
import { motion,  } from "framer-motion"
import { useState, useEffect,  } from "react"
import {  $defaultBlurCapRectangle, StyleSheetNode, motionDefaultProps } from "../../../../misc/types"
import { useResponsive } from "../../../../services/cssx/responsive/services/responsive/use-responsive"
import BlurFilter from "../../../../services/filter/blur/blur"
import NormalFilter from "../../../../services/filter/normal/normal"

interface input {
    src?: string,
    title?: string,
    alt?: string
    node?: motionDefaultProps<HTMLButtonElement>
    styleSheetNode?: StyleSheetNode
    delay?: number
    type?: "button" | "reset" | "submit"
}

export default function RectanglePostComposorButton({
    styleSheetNode, src, type, title, alt,
    node, delay }: input) {
    const [_src, setSrc] = useState<string | undefined>(src)
    const [_title, setTitle] = useState<string | undefined>(title)

    useEffect(() => {
        setSrc(src)
        setTitle(title)
    }, [src, title])


    const [click, setClick] = useState(false)
    const handleCombinedClick = (e: React.MouseEvent<HTMLButtonElement>) => {
        e.preventDefault()
        e.stopPropagation()

        setClick(true)

        const currentTarget = e.currentTarget
        const rect = currentTarget.getBoundingClientRect()

        if (node?.onClick) {
            const delegatedEvent = {
                ...e,
                currentTarget: {
                    ...currentTarget,
                    getBoundingClientRect: () => rect,
                },
            } as unknown as React.MouseEvent<HTMLButtonElement>

            node.onClick(delegatedEvent)
        }
    }
    const handleAnimationComplete = () => {
        setClick(false);
    }

    const [imgError, setImgError] = useState<boolean>(false)
    const { clamp, font, rectangle } = useResponsive()
    const tc = styleSheetNode?.btnTextColor ?? 'var(--foreground)'
    const th = styleSheetNode?.btnTheme ?? "var(--background)"
    const be = styleSheetNode?.blurEffect ?? false
    const bc = styleSheetNode?.blurCap ?? $defaultBlurCapRectangle

    const targetW = styleSheetNode?.btnWidth ?? 40
    const targetH = styleSheetNode?.btnHeight ?? 40
    const targetImg = styleSheetNode?.imgSize ?? 24
    const targetFont = styleSheetNode?.btnFontSize ?? 20
    const btnW = rectangle(targetW)
    const btnH = rectangle(targetH)
    const imgSize = clamp(targetImg, 'vh')
    const _fontSize = font(targetFont)

    const _tricolor = styleSheetNode?.triColor
    const { onClick: _onClick, ...buttonProps } = node ?? {}
    // ?? { right: colorPallete.$pink.floron_pink1, mid: colorPallete.$green.floron_green1, left: colorPallete.$green.floron_green2 }
    const r = clamp(styleSheetNode?.borderRadius??50, 'rem')
    return (
            <motion.button
                {...buttonProps}
                type={type ?? "button"}
                // whileHover={{ scale: 1.05 }}
                // whileTap={{ scale: .95 }}

                onClick={handleCombinedClick}
            style={{
                outline: 'none',
                border:'none',
                    width: btnW,
                    height: btnH,
                    borderRadius: r,
                backgroundColor: th,
                    backgroundImage:styleSheetNode?.imgTheme??'none'
                }}
                className="relative flex flex-row items-center p-2 justify-center gap-2  overflow-hidden"
            >
                {click && !be &&
                    <NormalFilter direction={styleSheetNode?.filterDirection} delay={delay} onComplete={handleAnimationComplete} tricolor={_tricolor} triggerAnimation={click} />
                }
                {click && be &&
                    <BlurFilter blurCap={bc} delay={delay} onComplete={handleAnimationComplete} tricolor={_tricolor} triggerAnimation={click} />
                }

                {(_src || _title) && (
                    <div className="absolute flex items-center justify-center z-20 pointer-events-none gap-2">
                        {_src && (
                            <div
                                className="rounded-full bg-white overflow-hidden flex-shrink-0"
                                style={{ width: imgSize, height: imgSize }}
                            >
                                {!imgError && (
                                    <img
                                        src={_src}
                                        alt={alt ?? "img"}
                                        width={80}
                                        height={80}
                                        className="w-full h-full object-cover"
                                        onError={() => setImgError(true)}
                                    />
                                )}
                            </div>
                        )}
                        {_title && (
                            <span
                                className="leading-none"
                                style={{ fontSize: _fontSize, color: tc }}
                            >
                                {_title}
                            </span>
                        )}
                    </div>
                )}
            </motion.button>
    )
}