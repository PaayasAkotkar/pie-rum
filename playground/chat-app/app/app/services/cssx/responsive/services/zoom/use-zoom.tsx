// @ copyright 2026 
import { useState, useEffect } from 'react';

// useZoomLevel this is created to have full control over the browser zoom
export function useZoomLevel() {
    const [zoom, setZoom] = useState(1);
    const [isMounted, setIsMounted] = useState(false);

    useEffect(() => {
        setIsMounted(true);
        if (typeof window === 'undefined') return;

        const updateZoom = () => {
            const dpr = window.devicePixelRatio;
            setZoom(dpr);
        };

        updateZoom();

        let mediaQuery: MediaQueryList | null = window.matchMedia(
            `(resolution: ${window.devicePixelRatio}dppx)`
        );

        const handleChange = () => {
            updateZoom();
            mediaQuery?.removeEventListener('change', handleChange);
            mediaQuery = window.matchMedia(
                `(resolution: ${window.devicePixelRatio}dppx)`
            );
            mediaQuery.addEventListener('change', handleChange);
        };

        mediaQuery.addEventListener('change', handleChange);
        window.addEventListener('resize', updateZoom);

        return () => {
            window.removeEventListener('resize', updateZoom);
            mediaQuery?.removeEventListener('change', handleChange);
        };
    }, []);

    return {
        zoom: isMounted ? zoom : 1,
        isMounted,
        counterScale: isMounted ? 1 / zoom : 1,
    };
}