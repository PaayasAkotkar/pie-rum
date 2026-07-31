// @ copyright 2026
export namespace CSSMaths {



    // GenerateClamp returns the clamp value based on params
    // targetPx: figma px or any as per your desired resp
    // baseDimension: display dim width 
    // baseZoom: clamping the value on based width of the display zoom
    // unit: any css unit prefered: -> rem, vh, vw, px
    // minRatio: resp zoom ratio
    // maxRatio: resp zoom ratio
    // baseDimensionH: display dm height
    // baseZoomH: same as baseZoom but for height
    // how things work is that you 
    export function GenerateClamp(
        targetPx: number,
        baseDimension: number,
        baseZoom: number,
        unit: string = 'vw',
        minRatio: number,
        maxRatio: number,
        baseDimensionH?: number,
        baseZoomH?: number,
    ): string {
        const isHeightUnit = unit === 'vh'

        const activeDimension = isHeightUnit ? (baseDimensionH ?? baseDimension) : baseDimension
        const activeZoom = isHeightUnit ? (baseZoomH ?? baseZoom) : baseZoom

        const virtualBase = activeDimension / activeZoom
        const dynamicV = (targetPx / virtualBase) * 100

        const min = (dynamicV * minRatio).toFixed(4)
        const mid = dynamicV.toFixed(4)
        const max = (dynamicV * maxRatio).toFixed(4)

        return `clamp(${min}${unit}, ${mid}${unit}, ${max}${unit})`
    }
    export const GridPattern = {
        DIAMOND_THIN: [1, 2],
        DIAMOND_WIDE: [2, 3],
        HEXAGON: [2, 3],
        TRIANGLE: [1, 2, 3, 4],
        INVERTED_TRIANGLE: [4, 3, 2, 1],
        HOURGLASS: [3, 2, 1, 2, 3],
        BRICK: [3],
    }
    export type GridPatternType = keyof typeof GridPattern;

    // GenerateGridPattern returns the 2d array of items to staggered
    // items: images to staggered
    // pattern: any geometric shape
    // @NOTE: make sure that to apply correct css in order to view it visually the same
    // @NOTE: use the 2d loop in the html template
    // @NOTE: you can also the provided GridPattern
    export function GenerateGridPattern<T>(items: T[], pattern: number[]): T[][] {
        const safePattern = pattern.length > 0 ? pattern : [1];

        const rows: T[][] = [];
        let i = 0;
        let step = 0;

        while (i < items.length) {
            let rowSize = safePattern[step % safePattern.length] || 1; // safety
            if (rowSize <= 0) rowSize = 1
            rows.push(items.slice(i, i + rowSize))
            i += rowSize
            step++
        }
        return rows;
    };

}


