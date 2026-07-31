// @ copyright 2026
'use client'
import { createContext, useContext, useCallback,  } from 'react'
import {
    mobileRepBaseHeight,
    mobileRepBaseWidth, mobileRespMax, mobileRespMin, mobileRespZoom,
    mobileRespZoomH,
    respBaseHeight,
    respBaseWidth, respMax, respMin, respZoom,
    respZoomH,
} from '../../algorithm/types'
import { useDevice } from '../device/use-device'
import { useZoomLevel } from '../zoom/use-zoom'
import { CSSMaths } from '../../algorithm/main'
export interface DeviceProfile {
    viewportW: number
    viewportH: number
    dpr: number
    isTouch: boolean
}

export interface SmartClampOptions {
    landscapePx: number
    portraitPx?: number      // defaults to landscapePx if omitted
    mobileLandscapePx?: number  // defaults to landscapePx
    mobilePortraitPx?: number   // defaults to mobileLandscapePx ?? landscapePx
}
// algorithm/main.ts  — new overload alongside GenerateClamp
export function GenerateSmartClamp(
    opts: SmartClampOptions,
    profile: DeviceProfile,
    baseDimension: number,
    baseZoom: number,
    unit: string = 'vw',
    minRatio: number,
    maxRatio: number,
    baseDimensionH?: number,
    baseZoomH?: number,
): string {
    const { viewportW, viewportH, isTouch } = profile
    const isLandscape = viewportW > viewportH
    const isMobile = isTouch && viewportW <= 1023

    // ── 1. pick the right design pixel value ──────────────────────────────────
    const lx = opts.landscapePx
    const ml = opts.mobileLandscapePx ?? lx
    const pot = opts.portraitPx ?? lx
    const mp = opts.mobilePortraitPx ?? ml

    const targetPx =
        isMobile
            ? isLandscape ? ml : mp
            : isLandscape ? lx : pot

    // ── 2. derive a device-aware correction scalar ─────────────────────────────
    //    Formula: compare actual viewport ratio to the "ideal" base ratio.
    //    This compensates for screens that are narrower/wider than the design base.
    const isHeightUnit = unit === 'vh'
    const activeDimension = isHeightUnit ? (baseDimensionH ?? baseDimension) : baseDimension
    //const activeZoom = isHeightUnit ? (baseZoomH ?? baseZoom) : baseZoom

    const idealRatio = activeDimension / (isHeightUnit ? baseDimension : (baseDimensionH ?? baseDimension))
    const actualRatio = isHeightUnit
        ? viewportH / viewportW
        : viewportW / viewportH

    // scalar nudges the target so the output *looks* correct on the real device
    // clamp scalar between 0.75–1.25 so it never goes wild
    const rawScalar = actualRatio / idealRatio
    const deviceScalar = Math.min(1.25, Math.max(0.75, rawScalar))

    const adjustedPx = targetPx * deviceScalar

    // ── 3. delegate to the existing GenerateClamp ─────────────────────────────
    return CSSMaths.GenerateClamp(
        adjustedPx,
        baseDimension,
        baseZoom,
        unit,
        minRatio,
        maxRatio,
        baseDimensionH,
        baseZoomH,
    )
}
// -----------------------------------------------------------------------------------------------------
// flow:
// place the ResponsiveProvider in the layout.tsx folder else nothing will work
// the rectangle, circle, square is provided to avoid using the resp stuff for buttons
// note: for circle,square you just use one single variable 
//       but again its all about prefering the way you like to use so go for it 
// font this is by far the best to avoid headache of fonts resp
// you can use w, h it still act same as clamp but i prefer you to use clamp directly incase 
// 
// -----------------------------------------------------------------------------------------------------

export type CSSUnit = 'vw' | 'vh' | 'px' | 'rem' | 'em' | '%'

export interface ClampOptions {
    base?: number
    baseZoom?: number
    minRatio?: number
    maxRatio?: number
}

export interface ResponsiveContextValue {
    zoom: number
    counterScale: number
    isMounted: boolean
    device: {
        isMobile: boolean
        isDesktop: boolean
        isTouch: boolean
        isLandscape: boolean
        isMobileLandscape: boolean
        isMobilePortrait: boolean
        isPcLandscape: boolean
        isPcPortrait: boolean
        is4K: boolean
        is2K: boolean
        isQHD: boolean
    }
    base: number
    baseZoom: number
    minRatio: number
    maxRatio: number
    clamp: (targetPx: number, unit?: CSSUnit, options?: ClampOptions) => string
    pick: (pc: number, mobileLandscape?: number, mobilePortrait?: number) => number
    font: (pc: number, mobileLandscape?: number, mobilePortrait?: number, unit?: CSSUnit) => string
    w: (targetPx: number, options?: ClampOptions) => string
    h: (targetPx: number, options?: ClampOptions) => string
    circle: (targetPx: number, options?: ClampOptions) => string
    rectangle: (targetPx: number, options?: ClampOptions) => string
    square: (targetPx: number, options?: ClampOptions) => string
    smart: (opts: SmartClampOptions, unit?: CSSUnit, overrides?: ClampOptions) => string
    smartFont: (opts: SmartClampOptions, unit?: CSSUnit) => string
    size: (targetPx: number) => { width: string; height: string }
}

const ResponsiveContext = createContext<ResponsiveContextValue | null>(null)
// -------------------------------------------------------------------------------------
// useResponsiveValues returns resp hand to generate the resp layout
// zoom,counterScale is kind of unorthdox and must only use if you know what you doing
// clamp: returns the clamp for target figma pixles
// pick: returns the dim
// font: returns the resp font [note: it uses vh or vw as per the layout]
// w,h returns width & height
// circle: returns the cricle shape
// rectangle returns the resp circle shape
// square: returns the resp square shape
// shorthand: for width and height -> use w,h instead
// -------------------------------------------------------------------------------------
function useResponsiveValues(): ResponsiveContextValue {
    const { zoom, counterScale, isMounted } = useZoomLevel()
    const { device } = useDevice()

    const isMobile = device.isMobile
    const mpot = device.isMobilePortrait
    const ppot = device.isPcPortrait
    const base = isMobile ? mobileRepBaseWidth : respBaseWidth
    const baseh = isMobile ? mobileRepBaseHeight : respBaseHeight
    const baseZoom = isMobile ? mobileRespZoom : respZoom
    const baseZoomH = isMobile ? mobileRespZoomH : respZoomH
    const minRatio = isMobile ? mobileRespMin : respMin
    const maxRatio = isMobile ? mobileRespMax : respMax
    const profile: DeviceProfile = {
        viewportW: typeof window !== 'undefined' ? window.innerWidth : base,
        viewportH: typeof window !== 'undefined' ? window.innerHeight : baseh,
        dpr: zoom,
        isTouch: device.isTouch,
    }

    const smart = useCallback(
        (opts: SmartClampOptions, unit: CSSUnit = 'vw', overrides: ClampOptions = {}): string => {
            if (!isMounted) return `${opts.landscapePx}${unit}`
            return GenerateSmartClamp(
                opts,
                profile,
                overrides.base ?? base,
                overrides.baseZoom ?? baseZoom,
                unit,
                overrides.minRatio ?? minRatio,
                overrides.maxRatio ?? maxRatio,
                baseh,
                baseZoomH,
            )
        },
        [isMounted, profile, base, baseh, baseZoom, baseZoomH, minRatio, maxRatio],
    )

    const smartFont = useCallback(
        (opts: SmartClampOptions, unit?: CSSUnit): string => {
            const autoUnit: CSSUnit = (device.isMobilePortrait || device.isPcPortrait) ? 'vw' : 'vh'
            return smart(opts, unit ?? autoUnit)
        },
        [smart, device],
    )
    const clamp = useCallback(
        (targetPx: number, unit: CSSUnit = 'vw', options: ClampOptions = {}): string => {
            if (!isMounted) return `${targetPx}${unit}`
            return CSSMaths.GenerateClamp(
                targetPx,
                options.base ?? base,
                options.baseZoom ?? baseZoom,
                unit,
                options.minRatio ?? minRatio,
                options.maxRatio ?? maxRatio,
                baseh,
                baseZoomH,
            )
        },
        [isMounted, base, baseh, baseZoom, baseZoomH, minRatio, maxRatio, smart, smartFont],
    )

    const pick = useCallback(
        (pc: number, mobileLandscape = pc, mobilePortrait = mobileLandscape): number => {
            if (device.isPcLandscape || device.isPcPortrait) return pc
            if (device.isMobileLandscape) return mobileLandscape
            return mobilePortrait
        },
        [device],
    )
    let unit: CSSUnit = mpot || ppot ? 'vw' : 'vh';

    const font = useCallback(
        (pc: number, mobileLandscape = pc, mobilePortrait = mobileLandscape, overrideUnit?: CSSUnit): string => {
            return clamp(pick(pc, mobileLandscape, mobilePortrait), overrideUnit ?? unit)
        },
        [clamp, pick, unit],
    )
    let unit2: CSSUnit = device.isLandscape ? 'vh' : 'vw';

    const circle = useCallback(
        (targetPx: number, options?: ClampOptions) => clamp(targetPx, unit2, options),
        [clamp, unit2],
    )

    let unit3: 'vw' | 'vh' = device.isLandscape ? 'vh' : 'vw'
    const rectangle = useCallback(
        (targetPx: number, options?: ClampOptions) => clamp(targetPx, unit3, options),
        [clamp, unit3],
    )

    let unit4: 'vw' | 'vh' = device.isLandscape ? 'vw' : 'vh'

    const square = useCallback(
        (targetPx: number, options?: ClampOptions) => clamp(targetPx, unit4, options),
        [clamp, unit4],
    )

    const w = useCallback(
        (targetPx: number, options?: ClampOptions) => clamp(targetPx, 'vw', options),
        [clamp],
    )

    const h = useCallback(
        (targetPx: number, options?: ClampOptions) => clamp(targetPx, 'vh', options),
        [clamp],
    )

    const size = useCallback(
        (targetPx: number) => ({ width: w(targetPx), height: h(targetPx) }),
        [w, h],
    )

    return {
        smart, smartFont,
        zoom, counterScale, isMounted,
        device, base, baseZoom, minRatio, maxRatio,
        clamp, pick, font, w, h, size, circle, rectangle, square
    }
}

export function ResponsiveProvider({ children }: { children: React.ReactNode }) {
    const value = useResponsiveValues()
    return (
        <ResponsiveContext.Provider value={value}>
            {children}
        </ResponsiveContext.Provider>
    )
}


export function useResponsive(): ResponsiveContextValue {

    const ctx = useContext(ResponsiveContext)
    if (!ctx) throw new Error('useResponsive must be used inside <ResponsiveProvider>')
    return ctx
}