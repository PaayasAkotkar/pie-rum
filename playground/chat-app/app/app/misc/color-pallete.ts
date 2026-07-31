export interface TriColor{
    right: string,
    mid: string,
    left:string
}
export namespace colorPallete {
    export const $family = {
        $a: {
            warmV1: "#F66E6E",
            warmV2: "#F6E26E",
            tri: {
                right: "#F66E70",
                left: "#F66EC4",
                mid: "#F6D26E"
            }
        },
        $b: {
            warmV1: "#A2F66E",
            warmV2: "#F6896E",
            tri: {
                right: "#6E89F6",
                left: "#6EB0F6",
                mid: "#6EF6B0"
            }
        },
        $c: {
            warmV1: "#6EB6F6",
            warmV2: "#A96EF6",
            tri: {
                right: "#AD6EF6",
                left: "#F6DD6E",
                mid: "#6EF670"
            }
        },
        $d: {
            warmV1: "#E4F66E",
            warmV2: "#6EF6ED",
            tri: {
                right: "#F66E6E",
                left: "#6EF697",
                mid: "#6E70F6"
            }
        },
        $e: {
            warmV1: "",
            warmv2: "",
            tri: {
                mid: "#FF7F00",
                left: "#2343FB",
                right: "#F2FF8F",

            }
        }
    }

    export const $orange = {
        v1: '#FFBE99',
        v2: '#FFB488',
        v3:"#FFAB7B",
    }

    export const $white = {
        v1: '#d2d2d2',
        v2: '#f0f0f0',
        v3:'#dbdada'
    }


    export const linearGradient = {
        galaxy: {
            bottom:
                `linear-gradient(to bottom, #DC2626, #F59E0B, #10B981, #059669, #DC2626)`,
                        top:
                `linear-gradient(to top, #DC2626, #F59E0B, #10B981, #059669, #DC2626)`,
            right:
                `linear-gradient(to right, #DC2626, #F59E0B, #10B981, #059669, #DC2626)`,

            left:
                `linear-gradient(to left, #DC2626, #F59E0B, #10B981, #059669, #DC2626)`,
            
    },
        
    }
}


