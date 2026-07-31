'use client'

import { colorPallete } from "../../../misc/color-pallete"
import { filter_direction } from "../../../misc/types"
import { useAnimate } from "framer-motion"
import { useEffect, useState } from "react"

export interface TriColor {
    right: string
    mid: string
    left: string
}
interface input {
    tricolor?: TriColor
    triggerAnimation: boolean
    onComplete?: () => void
    pause?: boolean
    blurCap?: string
    delay?: number
    direction?: filter_direction
}

// kept structurally identical (lead -> mid -> trail) across all directions
// so no direction feels "weaker" or "stronger" than another
const DIRECTION_ORDERS: Record<NonNullable<input["direction"]>, string[]> = {
    center: ["#mid", "#left", "#right"],
    leftToRight: ["#left", "#mid", "#right"],
    rightToLeft: ["#right", "#mid", "#left"],
    topToBottom: ["#mid", "#right", "#left"],
    bottomToTop: ["#mid", "#left", "#right"],
    radial: ["#mid", "#left", "#right"],
}

// travel offset (px) each blob enters FROM, per direction
const DIRECTION_OFFSETS: Record<NonNullable<input["direction"]>, { x: number, y: number }> = {
    center: { x: 0, y: 0 },
    leftToRight: { x: -120, y: 0 },
    rightToLeft: { x: 120, y: 0 },
    topToBottom: { x: 0, y: -120 },
    bottomToTop: { x: 0, y: 120 },
    radial: { x: 0, y: 0 },
}

// same stagger ratio applied uniformly to every direction
const STAGGER_RATIO = 0.18

export default function BlurFilter({ blurCap, pause, tricolor, triggerAnimation, onComplete, delay, direction = "center" }: input) {
    const _tricolor = tricolor || colorPallete.$family.$d.tri

    const [_pause, setPause] = useState(pause)
    const [isEntranceDone, setIsEntranceDone] = useState(false)
    const [scope, animate] = useAnimate()

    const order = DIRECTION_ORDERS[direction]
    const offset = DIRECTION_OFFSETS[direction]

    useEffect(() => {
        setPause(pause)
    }, [pause])

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
        setIsEntranceDone(true)
    }

    useEffect(() => {
        if (isEntranceDone && !_pause) {
            runExitAnimation().then(() => {
                setIsEntranceDone(false)
                if (onComplete) onComplete();
            })
        }
    }, [isEntranceDone, _pause])

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

    const blurCapacity = blurCap ?? "60px"

    return (
        <div ref={scope} className="w-full h-full z-1 pointer-events-none">
            <div className="absolute inset-0 w-full h-full">
                <div
                    id="right"
                    style={{
                        backgroundColor: _tricolor.right,
                        opacity: 0,
                        transform: `translate(${offset.x}px, ${offset.y}px) scale(0)`,
                        filter: `blur(${blurCapacity})`,
                        WebkitFilter: `blur(${blurCapacity})`,
                        willChange: 'transform, opacity',
                    }}
                    className="absolute top-0 left-0 w-[150%] h-[150%] rounded-full translate-x-[-20%] translate-y-[-10%]"
                />
                <div
                    id="mid"
                    style={{
                        backgroundColor: _tricolor.mid,
                        opacity: 0,
                        transform: `translate(${offset.x}px, ${offset.y}px) scale(0)`,
                        filter: `blur(${blurCapacity})`,
                        WebkitFilter: `blur(${blurCapacity})`,
                        willChange: 'transform, opacity',
                    }}
                    className="absolute top-0 left-0 w-[150%] h-[150%] rounded-full translate-x-[-45%] translate-y-[20%]"
                />
                <div
                    id="left"
                    style={{
                        backgroundColor: _tricolor.left,
                        opacity: 0,
                        transform: `translate(${offset.x}px, ${offset.y}px) scale(0)`,
                        filter: `blur(${blurCapacity})`,
                        WebkitFilter: `blur(${blurCapacity})`,
                        willChange: 'transform, opacity',
                        zIndex: 2,
                    }}
                    className="absolute bottom-0 right-0 w-[150%] h-[150%] rounded-full translate-x-[70%] translate-y-[15%]"
                />
            </div>
        </div>
    )
}