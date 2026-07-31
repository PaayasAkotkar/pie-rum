'use client'
import { StyleSheetNode } from "@/app/misc/types"
import { useResponsive } from "@/app/services/cssx/responsive/services/responsive/use-responsive"
import { useEffect, useState, useRef } from "react"


interface input {
    write: string
    startDelay?: number
    forceStop?: boolean
    onWritingStart?: () => void
    onWritingComplete?: () => void
    style?: StyleSheetNode
}

// Pencil returns a component text typewriter effect
// honestly speaking rxjs is quite better than NodeJS.Timeout as it saves time for implementation
// you can play around in rough how the typewriter works
// these are three ways to do so
// let text = "life goes on and one"
// let $i = 0

// const typwriter = setInterval(() => {
//     $i++;
//     let g = text.slice(0, $i);
//     console.log(g);

//     if ($i >= text.length) {
//         clearInterval(typwriter);
//     }
// }, 100);

// let stop = false
// interval(100).pipe(
//     takeWhile(() => !stop && $i < text.length),
//     tap(() => {
//         $i++
//         let g = text.slice(0, $i);
//         console.log(g);
//     })
// ).subscribe()

// from(text).pipe(
//     // concatMap is powerful in this case
//     // the reason being smiple not going to cancel and restart unlike switch map
//     // nor even going to start parallel like mergemap
//     concatMap(char => of(char).pipe(delay(100))),
//     // scane basically allows accumulation
//     scan((acc, curr) => acc + curr, "")
// ).subscribe(typewriter => {
//     console.log(typewriter);
// });
export default function Pencil({
    write,
    startDelay = 100,
    forceStop = false,
    onWritingStart,
    onWritingComplete,
    style
}: input) {
    const [text, setText] = useState<string>("")
    const _timeout = useRef<NodeJS.Timeout | null>(null)
    const $i = useRef<number>(0)
    const isWriting = useRef<boolean>(false)
    const previousWrite = useRef<string>("")
    const { font } = useResponsive()
    const f = font(style?.titleFontSize ?? 12)
    
    // stop animation
    useEffect(() => {
        if (forceStop && isWriting.current) {
            if (_timeout.current) {
                clearTimeout(_timeout.current)
                _timeout.current = null
            }
            isWriting.current = false
            onWritingComplete?.()
        }
    }, [forceStop])

    // start animation
    useEffect(() => {
        if (write !== previousWrite.current) {
            previousWrite.current = write
            setText("")
            $i.current = 0
            isWriting.current = false

            if (_timeout.current) {
                clearTimeout(_timeout.current)
                _timeout.current = null
            }
        }

        if (forceStop || !write || isWriting.current) {
            return
        }

        const startWriting = () => {
            isWriting.current = true
            onWritingStart?.()

            // alt:
            // 1. wait for startDelay
            // timer(startDelay).pipe(
            //     // 2. switch to the character stream
            //     concatMap(() => from(write)),
            //     // 3. add 100ms between each character
            //     concatMap(char => of(char).pipe(delay(100))),
            //     // 4. kill it if forceStop becomes true
            //     takeWhile(() => !forceStop),
            //     // 5. accumulate the string
            //     scan((acc, curr) => acc + curr, "")
            // ).subscribe({
            //     next: (text) => console.log(text),
            //     complete: () => console.log("Done!")
            // });

            const writeNextChar = () => {
                if ($i.current < write.length && !forceStop) {
                    $i.current++
                    setText(write.slice(0, $i.current))
                    // simple writing
                    _timeout.current = setTimeout(writeNextChar, 100)
                } else if ($i.current >= write.length) {
                    isWriting.current = false
                    onWritingComplete?.()
                }
            }
            // text write gap between delay
            _timeout.current = setTimeout(writeNextChar, 100)
        }

        // start delay
        _timeout.current = setTimeout(startWriting, startDelay)

        return () => {
            if (_timeout.current) {
                clearTimeout(_timeout.current)
                _timeout.current = null
            }
        }
    }, [write, forceStop])


    // imp for mobile devices else the color will be considered as white
    const textColor = style?.textColor ?? 'black'

    return (
        <div
            className="whitespace-pre-wrap"
            style={{ fontSize: f, color: textColor }}
        >
            {text}
        </div>
    )
}