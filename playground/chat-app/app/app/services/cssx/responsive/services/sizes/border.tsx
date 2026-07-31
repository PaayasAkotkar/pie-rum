
export interface BorderToken {
    radius: number
    borderWidth: number
    innerRadius: number
}


export const borders: BorderToken = {
    radius:10,
    borderWidth: 218,
    innerRadius: 188,
}

export function useBorders(): BorderToken {
    return borders
}
