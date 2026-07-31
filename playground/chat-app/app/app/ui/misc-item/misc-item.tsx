'use client'
import { colorPallete } from "@/app/misc/color-pallete"
import { useCopy } from "@/app/services/copy/copy"
import Pencil from "../pencil/pencil"
import MiscBackground from "./copy"
import { useState } from "react"
import { StyleSheetNode } from "@/app/misc/types"
import RectanglePostComposorButton from "../buttons/post-composor-buttons/rectangle-post-composor-button/rectangle-post-composor-button"
import { useResponsive } from "@/app/services/cssx/responsive/services/responsive/use-responsive"
import { ChessCoachMiscItems } from "@/app/services/rag-graphql/types"

interface input {
    misc: ChessCoachMiscItems
    startDelay?: number
    forceStop?: boolean
    style?: StyleSheetNode
}

export default function MiscItem({ style, startDelay, misc, forceStop = false }: input) {
    const { copy, isCopied } = useCopy()
    const _copy = isCopied ? 'copied' : 'copy'
    const [isCopyReady, setIsCopyReady] = useState(false)
    const [isLinkReady, setIsLinkReady] = useState(false)
    const { device,font } = useResponsive()
    const f=font(style?.titleFontSize??12)

    return (
        <div className="w-full max-w-full p-2">
            {misc.title && (
                <Pencil
                    startDelay={startDelay}
                    style={style}
                    write={misc.title}
                    forceStop={forceStop}
                />
            )}

            {misc.desc && (
                <Pencil
                    startDelay={startDelay}
                    style={style}
                    write={misc.desc}
                    forceStop={forceStop}
                />
            )}

            {misc.canCopy && misc.copy && (
                <MiscBackground>
                    <div
                        style={{
                            flexDirection: device.isMobilePortrait ? 'column' : 'row'
                        }}
                        className="flex gap-2 items-center">
                        <Pencil
                            startDelay={startDelay}
                            style={style}
                            write={misc.copy}
                            forceStop={forceStop}
                            onWritingComplete={() => setIsCopyReady(true)}
                        />
                        <RectanglePostComposorButton
                            styleSheetNode={style}
                            title={_copy}
                            node={{
                                onClick: () => copy(misc.copy ?? ""),
                                disabled: !isCopyReady,
                                style: { opacity: isCopyReady ? 1 : 0.5 }
                            }}
                        />
                    </div>
                </MiscBackground>
            )}

            {misc.isLink && misc.link && (
                <MiscBackground>
                    <a
                        href={isLinkReady ? misc.link : '#'}
                        target="_blank"
                        rel="noopener noreferrer"
                        onClick={(e) => !isLinkReady && e.preventDefault()}
                        style={{
                            fontSize: f,
                            color: isLinkReady ? colorPallete.$orange.v3 : 'gray',
                            cursor: isLinkReady ? 'pointer' : 'default'
                        }}
                        className="underline underline-offset-4 decoration-[1.1px]"
                    >
                        <Pencil
                            startDelay={startDelay}
                            style={style}
                            write={misc.link}
                            forceStop={forceStop}
                            onWritingComplete={() => setIsLinkReady(true)}
                        />
                    </a>
                </MiscBackground>
            )}
        </div>
    )
}
