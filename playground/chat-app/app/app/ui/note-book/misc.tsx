'use client'
import Pencil from "../pencil/pencil"
import MiscItem from "../misc-item/misc-item"
import { StyleSheetNode } from "@/app/misc/types"
import { ChessCoachReply } from "@/app/services/rag-graphql/types"

export default function Writer({
    note,
    startDelay,
    forceStop = false,
    onWritingStart,
    onWritingComplete,
    style
}: {
    note: ChessCoachReply | undefined,
    startDelay?: number,
    forceStop?: boolean,
    onWritingStart?: () => void,
    onWritingComplete?: () => void,
    style?: StyleSheetNode
}) {

    return (
        <>
            {note &&
                <div className="w-full">
                    <div className="whitespace-pre-wrap w-full break-words">
                        {note.year && (
                            <Pencil
                                startDelay={startDelay}
                                write={note.year}
                                forceStop={forceStop}
                                onWritingStart={onWritingStart}
                                onWritingComplete={onWritingComplete}
                                style={style}
                            />
                        )}
                        {note.title && (
                            <Pencil
                                startDelay={startDelay}
                                write={note.title}
                                forceStop={forceStop}
                                onWritingStart={onWritingStart}
                                onWritingComplete={onWritingComplete}
                                style={style}
                            />
                        )}
                        {note.desc && (
                            <Pencil
                                startDelay={startDelay}
                                write={note.desc}
                                forceStop={forceStop}
                                onWritingStart={onWritingStart}
                                onWritingComplete={onWritingComplete}
                                style={style}
                            />
                        )}
                        {note.miscItems && note.miscItems.map((item) => (
                            <MiscItem
                                startDelay={startDelay}
                                key={item.key}
                                misc={item.value}
                                forceStop={forceStop}
                                style={style}
                            />
                        ))}
                        {note.outro && (
                            <Pencil
                                startDelay={startDelay}
                                write={note.outro}
                                forceStop={forceStop}
                                onWritingStart={onWritingStart}
                                onWritingComplete={onWritingComplete}
                                style={style}
                            />
                        )}
                    </div>
                </div>
            }
        </>
    )
}