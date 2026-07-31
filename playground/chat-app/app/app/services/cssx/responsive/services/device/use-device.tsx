// @ copyright 2026
'use client'

import { useState, useEffect } from 'react';


export type DeviceShape = {
    isMobile: boolean
    isDesktop: boolean
    isTouch: boolean
    isLandscape: boolean
    isMobileLandscape: boolean
    isMobilePortrait: boolean
    isPcLandscape: boolean
    isPcPortrait: boolean
    is4K: boolean
    isQHD: boolean
    is2K: boolean
    isTvHD: boolean
}
// ─────────────────────────────────────────────────────────────────────────────
// Resolution breakpoints
//   Mobile portrait/landscape : width < 1024 (or touch + small)
//   PC Portrait               : 1024 ≤ width < 1440, desktop, portrait
//   PC Landscape              : 1440 ≤ width < 2300, desktop, landscape
//   TV HD                     : 2300 ≤ width < 2560  (covers 1920×1080 TVs
//                               rendered at device pixel ratio > 1, and
//                               actual 2K-ish panels at native res)
//   2K / QHD                  : 2560 ≤ width < 3840
//   4K                        : width ≥ 3840
//
// Why 2300 for isTvHD?
//   A 1920×1080 monitor reports innerWidth = 1920 at 100% zoom — that's
//   handled by isPcLandscape (1440–2299). Only screens wider than 2300
//   (e.g. a 1920-logical-px TV that exposes a wider viewport, or a native
//   2300+ panel) get the TV HD token. Pure UA sniffing for SmartTV strings
//   is kept as a secondary signal.
// ─────────────────────────────────────────────────────────────────────────────

const BREAKPOINTS = {
    mobileMax: 1023,       // ≤ 1023  → mobile
    pcPortraitMax: 1439,   // 1024–1439 desktop portrait
    pcLandscapeMax: 2299,  // 1440–2299 desktop landscape  ← 1920×1080 lives here
    tvHDMax: 2559,         // 2300–2559 TV HD
    qhdMax: 3839,          // 2560–3839 QHD / 2K
    // ≥ 3840 → 4K
} as const;

// useDevice returns true if any of the variable is matched by the dim
export function useDevice() {
    const [device, setDevice] = useState<DeviceShape>({
        isMobile: false,
        isDesktop: true,
        isTouch: false,
        isLandscape: true,
        isMobileLandscape: false,
        isMobilePortrait: false,
        isPcLandscape: true,
        isPcPortrait: false,
        is4K: false,
        is2K: false,
        isQHD: false,
        isTvHD: false,
    });

    useEffect(() => {
        const checkDevice = () => {
            const ua = typeof navigator === 'undefined' ? '' : navigator.userAgent;

            const isMobileUA = /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(ua);
            const isSmartTvUA = /TV|WebOS|Tizen|SmartTV|GoogleTV/i.test(ua);

            const isCoarsePointer = window.matchMedia('(pointer: coarse)').matches;
            const width = window.innerWidth;
            const height = window.innerHeight;
            const landscape = width > height;

            const isMobileDevice = isMobileUA || (isCoarsePointer && width <= BREAKPOINTS.mobileMax);
            const isDesktopDevice = !isMobileDevice;

            const is4K = width >= 3840;
            const isQHD = !is4K && width >= 2560;
            const is2K = isQHD;
            const isTvHD = !is4K && !isQHD && (isSmartTvUA || (!isMobileDevice && width >= 2300));

            const isRegularDesktop = isDesktopDevice && !is4K && !isQHD && !isTvHD;
            const isPcLandscape = isRegularDesktop && landscape && width >= 1440;
            const isPcPortrait = isRegularDesktop && (!landscape || width < 1440);

            setDevice({
                isMobile: isMobileDevice,
                isDesktop: isDesktopDevice,
                isTouch: isCoarsePointer,
                isLandscape: landscape,
                isMobileLandscape: isMobileDevice && landscape,
                isMobilePortrait: isMobileDevice && !landscape,
                isPcLandscape,
                isPcPortrait,
                is4K,
                isQHD,
                is2K,
                isTvHD,
            });

        };

        checkDevice();
        window.addEventListener('resize', checkDevice);
        return () => window.removeEventListener('resize', checkDevice);
    }, []);

    return { device };
}