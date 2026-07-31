import React from "react";
import { motion } from "framer-motion";

export interface IBackdropFilter {
    shadowColor?: string;
    shadowOpacity?: number
    spread?: string;
    blur?: string;
    theme?: string;
    direction?: "top" | "left" | "right" | "bottom" | "full";
}

interface BackdropFilterProps {
    children?: React.ReactNode;
    backdropFilter?: IBackdropFilter;
    radius?: string;
}
function applyIntensity(color: string | undefined, intensity: number): string {
    if (!color || color === "transparent") return "transparent";

    const alpha = Math.min(1, Math.max(0, intensity));

    // hex shorthand (#rgb or #rgba)
    const shortHex = color.match(/^#([0-9a-f]{3,4})$/i);
    if (shortHex) {
        const hex = shortHex[1];
        const r = parseInt(hex[0] + hex[0], 16);
        const g = parseInt(hex[1] + hex[1], 16);
        const b = parseInt(hex[2] + hex[2], 16);
        return `rgba(${r}, ${g}, ${b}, ${alpha})`;
    }

    // full hex (#rrggbb or #rrggbbaa)
    const fullHex = color.match(/^#([0-9a-f]{6,8})$/i);
    if (fullHex) {
        const hex = fullHex[1];
        const r = parseInt(hex.slice(0, 2), 16);
        const g = parseInt(hex.slice(2, 4), 16);
        const b = parseInt(hex.slice(4, 6), 16);
        return `rgba(${r}, ${g}, ${b}, ${alpha})`;
    }

    // rgb(r, g, b) → inject alpha
    const rgbMatch = color.match(/^rgb\(\s*(\d+),\s*(\d+),\s*(\d+)\s*\)$/i);
    if (rgbMatch) {
        return `rgba(${rgbMatch[1]}, ${rgbMatch[2]}, ${rgbMatch[3]}, ${alpha})`;
    }

    // rgba(r, g, b, a) → replace existing alpha
    const rgbaMatch = color.match(/^rgba\(\s*(\d+),\s*(\d+),\s*(\d+),\s*[\d.]+\s*\)$/i);
    if (rgbaMatch) {
        return `rgba(${rgbaMatch[1]}, ${rgbaMatch[2]}, ${rgbaMatch[3]}, ${alpha})`;
    }

    // named color or anything else — wrap in a way browsers handle via opacity fallback
    return color;
}
const directionStyles: Record<NonNullable<IBackdropFilter["direction"]>, React.CSSProperties> = {
    full: { inset: 0 },
    top: { top: 0, left: 0, right: 0, bottom: "auto", height: "40%" },
    bottom: { bottom: 0, left: 0, right: 0, top: "auto", height: "40%" },
    left: { left: 0, top: 0, bottom: 0, right: "auto", width: "40%" },
    right: { right: 0, top: 0, bottom: 0, left: "auto", width: "40%" },
};

export default function BackdropFilter({
    children,
    backdropFilter,
    radius,
}: BackdropFilterProps) {
    const inheritedRadius = radius
        ?? (React.isValidElement(children)
            ? (children.props as any)?.style?.borderRadius ?? "0px"
            : "0px");

    const b = backdropFilter?.blur ?? "20px";
    const s = backdropFilter?.spread ?? "0px";
    const d = backdropFilter?.direction ?? "full";
    const sc = backdropFilter?.shadowColor ?? "rgba(0, 0, 0, 0.5)";
    const shadowOpacity = backdropFilter?.shadowOpacity ?? 0.1;

    const resolvedShadowColor = (() => {
        const applied = applyIntensity(sc, shadowOpacity);
        if (applied === sc && shadowOpacity !== 1) {
            return `color-mix(in srgb, ${sc} ${Math.round(shadowOpacity * 100)}%, transparent)`;
        }
        return applied;
    })();

    const boxShadow = `0 0 ${b} ${s} ${resolvedShadowColor}`;
    return (
        <motion.div
            className="relative w-fit h-fit"
            initial="idle"
            whileHover="hovered"
        >
            <motion.div
                className="absolute pointer-events-none"
                variants={{
                    idle: { filter: "brightness(1)" },
                    hovered: { filter: "brightness(1.4)" },
                }}
                transition={{ duration: 0.2, ease: "easeOut" }}
                style={{
                    ...directionStyles[d],
                    position: "absolute",
                    backgroundColor: backdropFilter?.theme ?? "transparent",
                    boxShadow: boxShadow,
                    borderRadius: inheritedRadius,
                }}
            />
            <div className="relative z-10 w-full h-full">
                {children}
            </div>
        </motion.div>
    );
}