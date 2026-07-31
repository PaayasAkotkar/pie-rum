'use client'
import { useState, useEffect, } from "react"
import ProBookV1 from "../pro-books/pro-bookV1/pro-bookV1"
import { colorPallete } from "@/app/misc/color-pallete"
import MessageBubble from "../message-bubble/message-bubble"
import Writer from "./misc"
import { $defaultBorderRadius, StyleSheetNode } from "@/app/misc/types"
import { useResponsive } from "@/app/services/cssx/responsive/services/responsive/use-responsive"
import { OnChessCoachReply } from "@/app/services/rag-graphql/types"

interface InputProps {
    note?: OnChessCoachReply
    userNote: string
    onRead?: (value: string) => void
    writerStyle?: StyleSheetNode
    messageBubbleStyle?: StyleSheetNode
    frameStyle?: StyleSheetNode
}

interface IsClose {
    isClose: boolean
    Index: number
}

export default function NoteBook({ userNote, note, onRead, frameStyle, messageBubbleStyle, writerStyle }: InputProps) {
    const [history, setHistory] = useState<({ type: 'user', content: string } | { type: 'coach', content: OnChessCoachReply })[]>([])
    const [closedNotes, setClosedNotes] = useState<Set<number>>(new Set())
    const [stoppedWriters, setStoppedWriters] = useState<Set<number>>(new Set())
    const [isClosed, setIsClosed] = useState<Record<number, IsClose>>({})

    useEffect(() => {
        if (userNote) {
            // Stop all current writers when a new request is made
            setHistory(prev => {
                const newHistory = prev || [];
                // Add current indices to stoppedWriters
                setStoppedWriters(stopped => {
                    const next = new Set(stopped);
                    newHistory.forEach((_, i) => next.add(i));
                    return next;
                });
                return [...newHistory, { type: 'user', content: userNote }];
            });
        }
    }, [userNote])

    useEffect(() => {
        if (note) {
            setHistory(prev => {
                if (prev == null)
                    return [{ type: 'coach', content: note }]
                return [...prev, { type: 'coach', content: note }]
            })
        }
    }, [note])

    const handleToggleNote = (index: number) => {
        const isCurrentlyClosed = isClosed[index]?.isClose ?? false

        if (isCurrentlyClosed) {
            setIsClosed(prev => ({
                ...prev,
                [index]: { isClose: false, Index: index }
            }))

            setClosedNotes(prev => {
                const newSet = new Set(prev)
                newSet.delete(index)
                if (history[index + 1] && history[index + 1].type === 'coach') {
                    newSet.delete(index + 1)
                }
                return newSet
            })
        } else {
            // Force stop writing when closing
            const nextIndex = index + 1
            handleStopWriting(index)

            setTimeout(() => {
                setIsClosed(prev => ({
                    ...prev,
                    [index]: { isClose: true, Index: index }
                }))
                setClosedNotes(prev => new Set(prev).add(index))
                if (history[nextIndex] && history[nextIndex].type === 'coach') {
                    setClosedNotes(prev => new Set(prev).add(nextIndex))
                }
            }, 100)
        }
    }

    const handleStopWriting = (index: number) => {
        setStoppedWriters(prev => {
            const next = new Set(prev)
            next.add(index)
            // Also stop the coach reply which is index + 1
            next.add(index + 1)
            return next
        })
    }

    const { font, clamp } = useResponsive()
    const f = font(writerStyle?.titleFontSize ?? 12)
    const r = clamp($defaultBorderRadius)
    return (
        <ProBookV1
            style={frameStyle}
        >
            <div
                style={{
                    scrollbarColor: `black transparent`,
                    backgroundColor: frameStyle?.themeBg ?? 'transparent',
                    borderRadius: r,
                }}
                className="h-full w-full overflow-y-auto gap-2 p-2 overflow-x-hidden flex flex-col"
            >
                {history.map((item, index) => {
                    const isCurrentlyClosed = isClosed[index]?.isClose ?? false

                    if (item.type === 'user') {
                        return (
                            <div key={index} className="w-full">
                                {
                                    (item.content as string).length > 0 &&
                                    <MessageBubble
                                        style={messageBubbleStyle}
                                        textAreaStyle={frameStyle}
                                        query={isCurrentlyClosed ? `${(item.content as string).slice(0, 50)}...` : item.content as string}
                                        onEdit={(val) => onRead?.(val)}
                                        onClose={() => handleToggleNote(index)}
                                    />
                                }
                            </div>
                        )
                    } else {
                        const note = item.content as OnChessCoachReply
                        const shouldStop = stoppedWriters.has(index)
                        const isCoachClosed = closedNotes.has(index)
                        return (
                            <div key={index} className="w-full">
                                <div
                                    className="rounded-lg p-4 w-full max-w-full"
                                    style={{
                                        display: isCoachClosed ? 'none' : 'block',
                                        wordBreak: 'break-word',
                                        overflowWrap: 'break-word'
                                    }}
                                >
                                    {
                                        note.information &&
                                        <div className="flex flex-col gap-2 w-full">
                                            <Writer
                                                style={writerStyle}
                                                startDelay={50}
                                                note={note.information}
                                                forceStop={shouldStop}
                                            />
                                        </div>
                                    }
                                    {
                                        note.suggestion &&
                                        <div className="flex flex-col gap-2 w-full">
                                            <Writer
                                                style={writerStyle}
                                                startDelay={70}
                                                note={note.suggestion}
                                                forceStop={shouldStop}
                                            />
                                        </div>
                                    }
                                    {
                                        note.bestPractice &&
                                        <div className="flex flex-col gap-2 w-full">
                                            <Writer
                                                style={writerStyle}
                                                startDelay={80}
                                                note={note.bestPractice}
                                                forceStop={shouldStop}
                                            />
                                        </div>
                                    }
                                    {
                                        note.miscItems &&
                                        <div className="flex flex-col gap-2 w-full">
                                            <Writer
                                                style={writerStyle}
                                                startDelay={90}
                                                note={{ miscItems: note.miscItems }}
                                                forceStop={shouldStop}
                                            />
                                        </div>
                                    }
                                </div>
                                <div
                                    style={{ display: isCoachClosed ? 'block' : 'none' }}
                                    className="rounded-lg p-4 w-full cursor-pointer"
                                    onClick={() => handleToggleNote(index - 1)}
                                >

                                    <span style={{ fontSize: f }}>...</span>
                                </div>
                            </div>
                        )
                    }
                })}
            </div>
        </ProBookV1>
    )
}