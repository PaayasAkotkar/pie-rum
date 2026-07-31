'use client'

import { motion } from "framer-motion"
import { useState } from "react"
import { StyleSheetNode } from "@/app/misc/types"
import { useCopy } from "@/app/services/copy/copy"
import RectanglePostComposorButton from "../buttons/post-composor-buttons/rectangle-post-composor-button/rectangle-post-composor-button"
import { useResponsive } from "@/app/services/cssx/responsive/services/responsive/use-responsive"
import { colorPallete } from "@/app/misc/color-pallete"

interface input {
    query: string,
    onEdit?: (value: string) => void
    onClose?: () => void
    textAreaStyle?: StyleSheetNode
    style?: StyleSheetNode
}

export default function MessageBubble({
    textAreaStyle,
    style,
    onClose, query,
    onEdit }: input) {

    const [_query] = useState(query)
    const { font, clamp } = useResponsive()
    const tc = style?.textColor ?? 'black'
    const f = font(style?.titleFontSize ?? 12)
    const r = clamp(12)

    return (
        <div
                onClick={onClose}
            
            className="flex flex-col w-fit gap-2">
            <motion.div
                whileHover={{ filter: 'brightness(80%)' }}
                
                style={{
                    overflow: 'hidden',
                    borderRadius: r,
                    filter:'brightness(100%)',
                    backgroundColor: style?.theme ?? colorPallete.$orange.v1,
                }}
                className="w-fit p-3 bg-yellow1"
            >
                <p style={{
                    fontSize: f,
                    wordBreak: 'break-word', // Changed from break-all
                    whiteSpace: 'pre-wrap',
                    color: tc,
                }}>
                    {_query}
                </p>
            </motion.div>
        </div>
    )
}