export interface IConfigRequest {
    config: IConfig
    depth: number
    profileName: string
    kitName?: string
    serviceName?: string
    dispatcherName?: string
    eventName?: string
}
export interface IConfigResp {
    succeed: boolean
    error: string
}
export interface IConfig {
    activate: boolean
    swap?: ISwap
}

export interface ISwap {
    with: string
    current: string
    swap: boolean
}

export interface ChessStudentRequest {
  id: string
  name: string
  query: string;
}

export interface ChessCoachMiscItems {
  title?: string
  desc?: string
  canCopy?: boolean
  isLink?: boolean
}

export interface OnChessCoachReply {
  information?: ChessCoachReply
  suggestion?: ChessCoachReply
  bestPractice?: ChessCoachReply
  miscItems?: ChessCoachMiscEntry[]
}

export interface ChessCoachPayload {
  status: number
  message: string
}

export interface ChessCoachMiscEntry {
  key: string
  value: ChessCoachMiscItems
}

export interface ChessCoachMiscItems {
  title?: string
  desc?: string
  canCopy?: boolean
  isLink?: boolean
  link?: string
  copy?: string

}

export interface ChessCoachReply {
  year?: string
  title?: string
  desc?: string
  outro?: string
  miscItems?: ChessCoachMiscEntry[]
}
