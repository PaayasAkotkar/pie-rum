'use client'
import { ReactNode } from "react"
import { StyleSheetNode } from "../../misc/types";
import { useResponsive } from "@/app/services/cssx/responsive/services/responsive/use-responsive";

interface input {
    children: ReactNode
    styleSheetNode?: StyleSheetNode
    scroll?: boolean
    fitScreen?: boolean // New property to make it fill remaining flex space
}

export default function Frame({ children, scroll, styleSheetNode, fitScreen }: input) {

    const { clamp } = useResponsive()

    let baseW = styleSheetNode?.width ?? 593;
    let baseH = styleSheetNode?.height ?? 480;

    // If fitScreen is true, let flexbox handle sizing. Otherwise, calculate clamping.
    const w = fitScreen 
        ? '80vw' 
        : (styleSheetNode?.autoWidth == true ? '100%' : clamp(baseW, 'vw'))

    const h = fitScreen?'80vh': (styleSheetNode?.autoHeight == true ? '80vh' : clamp(baseH, 'vh'))

    return (
        <div
            style={{
                width: w,
                height: h,
                // flex: 1 shorthand tells the browser to grow and shrink to fit the leftover room
                flex: fitScreen ? '1 1 0%' : 'none', 
                overflow: scroll == true ? 'auto' : 'hidden',
            }}
        >
            {children}
        </div>
    )
}