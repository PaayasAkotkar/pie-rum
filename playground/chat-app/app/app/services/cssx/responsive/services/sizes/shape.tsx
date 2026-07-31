

export interface SizeVariant {
    width: number
    height: number
}

export interface ShapeSizeToken {
    xxxl: SizeVariant
    xxl: SizeVariant
    xl: SizeVariant
    large: SizeVariant
    medium: SizeVariant
    small: SizeVariant
    xs: SizeVariant
    xxs: SizeVariant
    xxxs: SizeVariant
}

export interface CircleSizeToken {
    xxxl: number
    xxl: number
    xl: number
    large: number
    medium: number
    small: number
    xs: number
    xxs: number
}

export interface AllShapeSizes {
    rectangle: ShapeSizeToken
    square: ShapeSizeToken
    circle: CircleSizeToken
}
export const shapes: AllShapeSizes = {
    rectangle: {
        xxxl: { width: 1860, height: 980 },   // near-full canvas, 1920 base
        xxl: { width: 1480, height: 780 },   // ×0.795 step
        xl: { width: 1160, height: 610 },   // ×0.784 step
        large: { width: 900, height: 475 },   // ×0.776 step
        medium: { width: 700, height: 370 },   // ×0.778 step
        small: { width: 540, height: 285 },   // ×0.771 step
        xs: { width: 400, height: 210 },   // ×0.741 step
        xxs: { width: 280, height: 148 },   // ×0.700 step
        xxxs: { width: 180, height: 95 },   // ×0.643 step — min readable
    },
    square: {
        xxxl: { width: 860, height: 860 },
        xxl: { width: 640, height: 640 },
        xl: { width: 480, height: 480 },
        large: { width: 360, height: 360 },
        medium: { width: 270, height: 270 },
        small: { width: 200, height: 200 },
        xs: { width: 150, height: 150 },
        xxs: { width: 100, height: 100 },
        xxxs: { width: 60, height: 60 },
    },
    circle: {
        xxxl: 860,
        xxl: 640,
        xl: 480,
        large: 360,
        medium: 270,
        small: 200,
        xs: 150,
        xxs: 100,
    },
}
export function useShapeSizes(): AllShapeSizes {
    return shapes
}