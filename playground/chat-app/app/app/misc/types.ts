import React from "react";
import { HTMLMotionProps } from "framer-motion"
export type page = string

export type IMenu = {
    title: string,
    sheet?: StyleSheetNode
}
export type PostDirection = 'top-left' | 'top-right' | 'bottom-left' | 'bottom-right'
export interface IError {
    message: string

}
export interface IPost {
    title: string
    subTitle: string
    drection: PostDirection
    btnTitle: string
    label: string
    origin?: string
    onPost?: () => void
    link?: ILink
    styles?: StyleSheetNode
}
export interface IWrite {
    text: string
    underline?: boolean
    bold?: boolean
    uppercase?: boolean
    isTitle?: boolean
    prefix?: string
    breakText?: boolean
    breakTextDirection?: 'inline' | 'block'
    breakIndex?: number
    textAlign?: 'left' | 'right' | 'center'
    styles?: StyleSheetNode
}
export type motionDefaultProps<T extends HTMLElement> = HTMLMotionProps<any> & React.HTMLAttributes<T>
export interface ILink {
    href: string
    content: string
    newTab?: boolean
    mailTo?: boolean
}
export interface IFooter {
    title: string
    content: ILink[]
}
export type defaultProps<T extends HTMLElement> = React.HTMLAttributes<T>


export interface IContent {
    title: string
    desc: string[]
}
export interface ISession {
    pickToss: boolean,
    toss: boolean,
    dictionary: boolean,
    challenge: boolean,
    guessStage: boolean
}
export interface IAdvertise {
    src: string
    alt: string
    link?: string
    title: string
    desc: string[]
    tag?: string
}
export interface StyleSheetNode {
    blurEffect?: boolean
    gap?: number
    imgTheme?: string
    btnWidth?: number
    btnHeight?: number
    width?: number
    height?: number
    fontFamily?: string
    direction?: "top" | "bottom" | "left" | "right"
    hoverTheme?: string
    tapTheme?: string
    miscWidth?: number
    padding?: number
    miscHeight?: number
    lineHeight?: number
    borderRadius?: number
    autoHeight?: boolean
    autoWidth?: boolean
    frameSize?: number
    size?: number // this can be use for the button size 
    imgSize?: number
    imgWidth?: number
    imgheight?: number

    // typography
    titleFontSize?: number
    subTitleFontSize?: number
    labelFontSize?: number
    textColor?: string
    btnTheme?: string
    // end
    btnFontSize?: number

    scrollBarColor?: string
    theme?: string
    themeBg?: string // used for another backdiv
    blurCap?: string
    btnTextColor?: string
    textFrameTheme?: string
    triColor?: { right: string, left: string, mid: string }
    boxShadowColor?: string
    activeTextColor?: string
    inactiveTextColor?: string
    activeColor?: string
    inactiveColor?: string
    borderColor?: string
    borderWidth?: number
    filterDirection?: filter_direction
}
export type filter_direction = "center" | "leftToRight" | "rightToLeft" | "topToBottom" | "bottomToTop" | "radial"

export interface IInvite {
    title: string
    desc: string
}
export interface IPolicy {
    policy: string[],
    termsAndServices: string[],
    FAQ: string[],
    feedbackAndComplaints: string[]
}
export interface INewsLetter {
    date: IWrite
    title: IWrite
    author: IWrite
    story: IWrite[]
    imgs?: string[]
}

export interface IArticle {
    title: IWrite
    subTitle: string
    desc: IWrite
    topColor?: string
    bottomColor?: string
}

// for overlay
export const $glassEffect: { color: string } = {
    color: "rgba(130, 163, 242, 0.3)"
}
export interface IForm {
    salary: string
    the_role: string
    location: string

    key_responsibilities: string[]
    what_we_offer: string[]
    culture_and_vibes: string
}

export interface IFormRequest {
    topic: string
    fields: IForm
    isCustom: boolean
}

export interface effect {
    blur?: number,
    color?: string
    saturation?: number
}

export interface chip {
    title: IWrite
}

export interface Dispatcher {
    play: () => void
    event: keyboard
    // any combo or single
    // combo=> control + C
    // single=> enter
    combo: string
    // end
    uiName: string // if used for UI what will be its name
    key: string // name of the current Dispatcher
}

type keyboard = "keydown" | "keypress" | "keypress"

export interface DispatcherDocument {
    name: string
    combo: string // if any
}

export const $delay = 1200
export const $defaultBorderWidth = 4
export const $defaultBlurCap = '10px'
export const $defaultBlurCapRectangle = '20px'
export const $defaultBorderRadius = 15
export const $defaultDescTitle = 30
export const $defaultDescDesc = 25
export const $squareBorderRadius = 10


export const $muteKey = 'mute'
export const $unmuteKey = 'unmute'
export const $activeSpeaking = 'active-speaking'
export const $transcripts = 'transcripts'
export const $lobbyStatus = 'transcripts'
export const $idKey = "session-login-id"
export const $defaultIDLen = 5

export interface Void {
    info: FileSliceInfo,
}

export interface FileSliceInfo {
    sys_info: SysInfo,
    slice_info: SliceInfo,
}

export interface SliceInfo {
    ext: string,             // extension of the file we reading into & splitting
    mark_ext: string,        // extension that is wrote after the ext
    slices: Slice[],      // actually parts of the file
    stored_location: string, // location where the file parts are stored
    hash: string,
}

export interface FileSliceInfo {
    sys_info: SysInfo,
    slice_info: SliceInfo,
}

export interface Slice {
    data: any,
}

export interface SysInfo {
    space_left: number,
    total_space: number,
    drive: string,
    src: string, // src-from which the requesst made
}