'use client'
import { ReactNode, useEffect, useRef, useState } from "react";
import { SubmitHandler, useForm } from "react-hook-form";
import { useResponsive } from "@/app/services/cssx/responsive/services/responsive/use-responsive";
import { $defaultBorderRadius, StyleSheetNode } from "@/app/misc/types";
import CirclePostComposorButton from "../buttons/post-composor-buttons/circle-post-composor-button/circle-post-composor-button";
import Frame from "../frame/frame";

interface input {
    get?: (value: string) => void
    defaultValue?: string
    styleSheetNode?: StyleSheetNode
    minLength?: number
    maxLength?: number
    placeHolder?: string
    onCancel?: () => void
    title?: string
    row?: number
    children?:ReactNode
}

export default function TextAreaV1({children, get, row, title, defaultValue, styleSheetNode, minLength, maxLength, placeHolder, onCancel }: input) {

    const { register, reset, handleSubmit } = useForm<{ text: string }>()
    const textareaRef = useRef<HTMLTextAreaElement | null>(null)
    const [ph, setPh] = useState(placeHolder ?? "go...")

    useEffect(() => {
        setPh(placeHolder ?? "go...")
    }, [placeHolder])

    const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
        if (e.key === 'Enter' && !e.shiftKey) {
            e.preventDefault()
            handleSubmit(onSubmit)()
        }
    }

    const onSubmit: SubmitHandler<{ text: string }> = (data: { text: string }) => {
        get?.(data.text)
        reset({ text: "" })
    }

    const eraseText = () => {
        reset({ text: "" })
        onCancel?.()
    }
    useEffect(() => {
        const eraseEvent = (ev: KeyboardEvent) => {
            if (ev.key.toLowerCase() == 'delete') {
                ev.preventDefault()
                eraseText()
            }
        }

        window.addEventListener('keydown', eraseEvent)
        return () => {
            window.removeEventListener("keydown", eraseEvent);
        }
    }, [])

    const { font, clamp } = useResponsive()
    const fontSize = font(styleSheetNode?.titleFontSize ?? 40)
    const textColor = styleSheetNode?.textColor ?? 'var(--foreground)'
    const r = clamp($defaultBorderRadius)
    return (
        <Frame styleSheetNode={{ ...styleSheetNode, frameSize: styleSheetNode?.frameSize ?? 500, }}>
            <div style={{
                borderRadius: r,
                backgroundColor: styleSheetNode?.theme ?? 'var(--primary)'
            }} className="w-full  h-full p-1 flex flex-col" >
                <textarea
                    {...register('text')}
                    ref={(e) => {
                        register('text').ref(e)
                        textareaRef.current = e
                    }}
                    defaultValue={defaultValue ?? ""}
                    onKeyDown={(e) => handleKeyDown(e)}
                    rows={row ?? 1}
                    autoFocus={true}
                    minLength={minLength}
                    maxLength={maxLength}
                    placeholder={ph}
                    style={{
                        scrollbarColor: 'black transparent',
                        fontSize: fontSize,
                        color: textColor,
                        //height: styleSheetNode?.height ?? 200,
                        backgroundColor: 'transparent',
                        overflowY: 'auto',
                    }}
                    className='resize-none w-full h-full focus:outline-none pl-3'
                />
                <div className='w-full h-fit relative right-1 p-2 flex flex-row justify-between'>
                    <div>
                        <CirclePostComposorButton
                            type="button"
                            styleSheetNode={{...styleSheetNode }}
                            title="rem"
                            node={{ onClick: eraseText }}
                        />
                    </div>
                    {
                        children&&children
                    }
                    <div>
                        <CirclePostComposorButton
                            type="submit"
                            styleSheetNode={{  ...styleSheetNode }}
                            title={title ?? "don"}
                            node={{ onClick: handleSubmit(onSubmit) }}
                        />
                    </div>
                </div>
            </div>
        </Frame>
    )
}