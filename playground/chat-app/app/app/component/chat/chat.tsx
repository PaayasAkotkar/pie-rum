'use client'
import { colorPallete } from "@/app/misc/color-pallete";
import { StyleSheetNode } from "@/app/misc/types";
import { useResponsive } from "@/app/services/cssx/responsive/services/responsive/use-responsive";
import { useRag, useRagSubscription } from "@/app/services/rag-graphql/rag-graphql";
import { OnChessCoachReply } from "@/app/services/rag-graphql/types";
import CirclePostComposorButton from "@/app/ui/buttons/post-composor-buttons/circle-post-composor-button/circle-post-composor-button";
import NoteBook from "@/app/ui/note-book/note-book";
import TextArea from "@/app/ui/text-area/text-area";
import { useEffect, useState } from "react";

// Chat returns the chat system with ai
export default function Chat() {

    let $name = "Pie_Rum_Fella@26"

    let $id = "xxx222"

    const [note, setNote] = useState('')

    const [coachNote, setCoachNote] = useState<OnChessCoachReply | undefined>()

    const { mut } = useRag()

    const { aiReply } = useRagSubscription({
        id: $id,
        name: $name,
        query: note,
    })

    useEffect(() => {
        setCoachNote(aiReply)
    }, [aiReply])

    const handleSend = async (text: string) => {
        setNote(text)
        try {
            const result = await mut.askAI({
                id: $id,
                query: text,
                name: $name
            });
            console.log("Message sent successfully:", result);
        } catch (error) {
            console.error("Error sending message:", error);
        }
    }

    // just pass this to test
    // @ts-ignore
    let testNote1: OnChessCoachReply = {
        information: {
            desc: "Hello Pie_Rum_Fella! I've put together a list of some of the best YouTube channels for chess content that I think you'll find incredibly helpful for improving your game and staying entertained. These creators offer a wide range of styles, from beginner-friendly guides to in-depth master game analysis."
        },
        miscItems: [
            {
                key: "yt_gothamchess_channel",
                value: {
                    title: "📺 GothamChess Channel",
                    desc: "Levy Rozman (GothamChess) is known for his engaging and often humorous explanations, making complex concepts accessible. Great for all levels, especially beginners and intermediate players looking to improve in 2025.",
                    canCopy: false,
                    isLink: true,
                    link: "https://www.youtube.com/@GothamChess"
                }
            },
            {
                key: "yt_agadmator_channel",
                value: {
                    title: "📺 Agadmator's Chess Channel",
                    desc: "Ante Saric (Agadmator) provides calm, insightful analysis of famous games, current events, and historical matches. Perfect for those who enjoy in-depth game breakdowns and learning about chess history in 2025.",
                    canCopy: false,
                    isLink: true,
                    link: "https://www.youtube.com/@agadmator"
                }
            },
            {
                key: "yt_naroditsky_channel",
                value: {
                    title: "📺 Daniel Naroditsky (Danya) Channel",
                    desc: "Grandmaster Daniel Naroditsky offers high-level instructive content, including speedruns, game analysis, and detailed opening theory. Excellent for serious improvement and strategic understanding in 2025.",
                    canCopy: false,
                    isLink: true,
                    link: "https://www.youtube.com/@DanielNaroditsky"
                }
            },
            {
                key: "yt_ericrosen_channel",
                value: {
                    title: "📺 Eric Rosen Channel",
                    desc: "Eric Rosen combines fun, educational content with brilliant tactical puzzles and interesting game commentary. His 'Saint Louis Chess Club' content is also top-tier for 2025.",
                    canCopy: false,
                    isLink: true,
                    link: "https://www.youtube.com/@EricRosen"
                }
            },
            {
                key: "yt_video_opening_principles",
                value: {
                    title: "📺 Video: The ULTIMATE Guide To Opening Principles",
                    desc: "This video from GothamChess breaks down the fundamental principles of chess openings in an easy-to-understand way, crucial for building a strong foundation for your games in 2025.",
                    canCopy: false,
                    isLink: true,
                    link: "https://www.youtube.com/watch?v=H7oB1tN3K4A"
                }
            }
        ]
    }

    // @ts-ignore
    const testNote2: OnChessCoachReply = {
        information: {
            title: "Welcome, Coach is here!",
            desc: "Hello Pie_Rum_Fella! Yes, I'm here and ready to help you with your chess journey. What's on your mind today, or what would you like to work on? Let's make 2025 your best chess year yet!"
        },
        miscItems: [
            {
                key: "youtube_fundamentals",
                value: {
                    title: "📺 Chess Fundamentals for Beginners",
                    desc: "This video provides a great overview of essential chess concepts every player should know.",
                    canCopy: false,
                    isLink: true,
                    link: "https://www.youtube.com/watch?v=kCO-92kQ6-o"
                }
            },
            {
                key: "book_mastering_chess",
                value: {
                    title: "📚 Mastering Chess: A Complete Guide",
                    desc: "A comprehensive book covering strategy, tactics, and endgames for aspiring players.",
                    canCopy: false,
                    isLink: false
                }
            },
            {
                key: "fen_sample_middlegame",
                value: {
                    title: "📋 FEN for Middlegame Analysis",
                    desc: "Copy this FEN into your favorite analysis board (like Lichess or Chess.com) to practice your tactical vision in a complex middlegame position.",
                    canCopy: true,
                    isLink: false,
                    copy: "r1bqkb1r/pppn1ppp/3p1n2/4p3/3P4/2NBPN2/PPP2PPP/R1BQK2R b KQkq - 0 5"
                }
            },
            {
                key: "lichess_sicilian_defense",
                value: {
                    title: "🔗 Explore the Sicilian Defense (Lichess)",
                    desc: "Dive into the Sicilian Defense on Lichess's opening explorer to see common variations and master lines.",
                    canCopy: false,
                    isLink: true,
                    link: "https://lichess.org/opening/Sicilian_Defense"
                }
            }
        ]
    }

    // @ts-ignore
    const testUserNote = 'pass me some of the best content creators links on youtube.'
    // end

    const { clamp } = useResponsive()

    const f = clamp(60)

    const frameSheet: StyleSheetNode = {
        width: 1400,
        height: 600,
        theme: colorPallete.$white.v1,
        themeBg: colorPallete.$white.v2,
    }

    const blurCap = "12px"

    const bubbleMessageSheet: StyleSheetNode = {
        theme: colorPallete.$white.v3,
        blurEffect: true,
        blurCap: blurCap,
        btnWidth: 100,
        titleFontSize: 50,
        size: 50,
        btnTextColor: 'white'
    }

    const textAreaSheet: StyleSheetNode = {
        width: 800,
        height: 240,
        blurEffect: true,
        blurCap: blurCap,
        frameSize: 30,
        titleFontSize: 49,
        size: 80, btnFontSize: 40,
        theme: colorPallete.$white.v2,
        btnTextColor: 'white',
        textColor: 'black'
    }

    const ws: StyleSheetNode = {
        titleFontSize: 45,
        btnWidth: 140,
        btnHeight: 50,
        btnFontSize: 40,
    }
    const handleConfig = () => {
        mut.updateConfig({
            config: {
                activate: true,
                swap: {
                    swap: false,
                    with: "",
                    current:"",
                }
            },
            depth: 2,
            profileName: "nvidea",
            kitName: "nvidea-kit",
            serviceName: "nvidea-service",
            dispatcherName: "nvidea-dispatcher",
            eventName: "nvidea-event",
        },
        )
    }
    return (
        <>
            <div className="inset-0 h-screen flex flex-col justify-center items-center p-1 gap-5">
                <span style={{
                    fontWeight: 'bold',
                    color: 'black',
                    fontSize: f
                }}>PIE RUM CHESS TALKS</span>

                <NoteBook
                    frameStyle={frameSheet}
                    messageBubbleStyle={bubbleMessageSheet}
                    writerStyle={ws}
                    note={coachNote}
                    userNote={note}
                    onRead={handleSend}
                ></NoteBook>

                <div className="w-full flex flex-col justify-center items-center">
                    <TextArea
                        title="ask"
                        get={handleSend}
                        styleSheetNode={textAreaSheet}
                    >

                        <CirclePostComposorButton node={{
                            onClick: () => {
                                handleConfig()
                            }
                        }} styleSheetNode={textAreaSheet} title="con" />
                    </TextArea>
                </div>
            </div>
        </>
    )
}
