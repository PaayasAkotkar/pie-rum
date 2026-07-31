'use client'

import { colorPallete } from "../../../misc/color-pallete"
import { filter_direction } from "../../../misc/types"
import { TriColor } from "../../../misc/color-pallete"
import { useAnimate } from "framer-motion"
import { useEffect } from "react"

interface input {
    tricolor?: TriColor,
    triggerAnimation: boolean,
    onComplete?: () => void
    delay?: number
    direction?: filter_direction
}

// order in which the three blobs animate, kept structurally identical
// (always mid -> secondary -> tertiary) so stagger weight feels equal
const DIRECTION_ORDERS: Record<NonNullable<input["direction"]>, string[]> = {
    center: ["#mid", "#left", "#right"],
    leftToRight: ["#left", "#mid", "#right"],
    rightToLeft: ["#right", "#mid", "#left"],
    topToBottom: ["#mid", "#right", "#left"],
    bottomToTop: ["#mid", "#left", "#right"],
    radial: ["#mid", "#left", "#right"],
}

// travel offset (in px) each blob enters FROM, per direction.
// center/radial = no travel, just scale+fade in place.
const DIRECTION_OFFSETS: Record<NonNullable<input["direction"]>, { x: number, y: number }> = {
    center: { x: 0, y: 0 },
    leftToRight: { x: -120, y: 0 },
    rightToLeft: { x: 120, y: 0 },
    topToBottom: { x: 0, y: -120 },
    bottomToTop: { x: 0, y: 120 },
    radial: { x: 0, y: 0 },
}

// every direction uses the SAME stagger gap ratio and spring config
// so none feels weaker/stronger than another
const STAGGER_RATIO = 0.18 // 18% of entrance/exit time, applied uniformly

export default function NormalFilter({ tricolor, triggerAnimation, onComplete, delay, direction = "center" }: input) {
    const _tricolor = tricolor || colorPallete.$family.$d.tri
    const [scope, animate] = useAnimate()

    const order = DIRECTION_ORDERS[direction]
    const offset = DIRECTION_OFFSETS[direction]

    const runEntranceAnimation = async () => {
        const total = delay ?? 1000
        const tEnt = total * 0.6
        const gap = (tEnt * STAGGER_RATIO) / 1000

        await Promise.all(
            order.map((id, i) =>
                animate(
                    id,
                    { scale: 1, opacity: 1, x: 0, y: 0 },
                    {
                        type: "spring",
                        stiffness: 220,
                        damping: 18,
                        mass: 0.7,
                        delay: gap * i,
                    }
                )
            )
        )

        await runExitAnimation()
        if (onComplete) onComplete()
    }

    const runExitAnimation = async () => {
        const total = delay ?? 1000
        const tExt = total * 0.4
        const gap = (tExt * STAGGER_RATIO) / 1000
        const exitOrder = [...order].reverse()

        await Promise.all(
            exitOrder.map((id, i) =>
                animate(
                    id,
                    { scale: 0, opacity: 0, x: offset.x, y: offset.y },
                    {
                        type: "spring",
                        stiffness: 320,
                        damping: 22,
                        mass: 0.6,
                        delay: gap * i,
                    }
                )
            )
        )
    }

    useEffect(() => {
        if (triggerAnimation) {
            runEntranceAnimation()
        }
    }, [triggerAnimation, direction])

    return (
        <div ref={scope} className="absolute inset-0 w-full h-full z-1 pointer-events-none">
            <div
                id="right"
                style={{
                    backgroundColor: _tricolor.right,
                    opacity: 0,
                    transform: `translate(${offset.x}px, ${offset.y}px) scale(0)`,
                    willChange: 'transform, opacity'
                }}
                className="absolute top-0 left-0 w-[150%] h-[150%] rounded-full translate-x-[-20%] translate-y-[-10%]" />

            <div
                id="mid"
                style={{
                    backgroundColor: _tricolor.mid,
                    opacity: 0,
                    transform: `translate(${offset.x}px, ${offset.y}px) scale(0)`,
                    willChange: 'transform, opacity'
                }}
                className="absolute top-0 left-0 w-[150%] h-[150%] rounded-full translate-x-[-55%] translate-y-[20%] z-2" />

            <div
                id="left"
                style={{
                    backgroundColor: _tricolor.left,
                    opacity: 0,
                    transform: `translate(${offset.x}px, ${offset.y}px) scale(0)`,
                    willChange: 'transform, opacity'
                }}
                className="absolute bottom-0 right-0 w-[150%] h-[150%] rounded-full translate-x-[70%] translate-y-[15%]" />
        </div>
    )
}