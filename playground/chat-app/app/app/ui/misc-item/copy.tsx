'use client'

import { useResponsive } from "@/app/services/cssx/responsive/services/responsive/use-responsive"


interface CopyBackgroundProps {
    children: React.ReactNode
}

export default function MiscBackground({ children }: CopyBackgroundProps) {
    const { clamp,device } = useResponsive()

    const roundedValue = device.isPcLandscape ? 50 : device.isMobilePortrait ? 10 : 12
    const rounded = clamp(roundedValue)

    const paddingValue = device.isPcLandscape ? 4 : device.isMobilePortrait ? 16 : 20
    const padding = clamp(paddingValue)

    return (
        <div
            className="w-fit bg-yellow1 overflow-hidden"
            style={{
                borderRadius: rounded,
                padding: padding
            }}
        >
            {children}
        </div>
    )
}