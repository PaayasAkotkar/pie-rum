'use client'
import { colorPallete } from "@/app/misc/color-pallete"
import { useRef, } from "react"
import { $defaultBorderRadius, StyleSheetNode } from "@/app/misc/types"
import { useResponsive } from "@/app/services/cssx/responsive/services/responsive/use-responsive"

interface input {
    style?: StyleSheetNode
    children?: React.ReactNode
}

export default function ProBookV1({  children, style }: input) {
    const { clamp } = useResponsive()
    
    const w =clamp(style?.width??12)
    const h = clamp(style?.height ?? 10)
    
    const rounded = clamp($defaultBorderRadius)

    const _padding = clamp(20)

    return (
        <>
            <div style={{
                width: w,
                height: h,
                backgroundColor: style?.theme?? colorPallete.$orange.v1,
                borderRadius: rounded,
                padding: _padding,
                overflow: 'hidden'
            }}>
                    <div className="w-full h-full relative flex flex-col justify-between overflow-hidden rounded-b-[inherit]">
                        <div className="relative w-full flex-1 min-h-0 flex flex-col">
                            {children}
                        </div>
                </div>
            </div>

        </>
    )
}