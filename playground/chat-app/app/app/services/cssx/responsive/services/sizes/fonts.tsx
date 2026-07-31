export interface ScreenTitleVariant {
    titleFontSize: number
    subTitleFontSize: number
}

export interface ScreenTitleToken {
    xxxl: ScreenTitleVariant
    xxl: ScreenTitleVariant
    xl: ScreenTitleVariant
    large: ScreenTitleVariant
    medium: ScreenTitleVariant
    small: ScreenTitleVariant
    xs: ScreenTitleVariant
    xxs: ScreenTitleVariant
}
export const screenTitles: ScreenTitleToken = {
    xxxl: { titleFontSize: 334, subTitleFontSize: 220 },
    xxl: { titleFontSize: 200, subTitleFontSize: 128 },
    xl: { titleFontSize: 158, subTitleFontSize: 99 },
    large: { titleFontSize: 112, subTitleFontSize: 70 },
    medium: { titleFontSize: 78, subTitleFontSize: 50 },
    small: { titleFontSize: 58, subTitleFontSize: 35 },
    xs: { titleFontSize: 43, subTitleFontSize: 27 },
    xxs: { titleFontSize: 29, subTitleFontSize: 18 },
}

export function useScreenTitleSizes(): ScreenTitleToken {
    return screenTitles
}