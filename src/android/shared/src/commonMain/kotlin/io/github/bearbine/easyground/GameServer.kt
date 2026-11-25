package io.github.bearbine.easyground.model

// игровой сервер, размещенный на хосте
data class GameServer(
    val id:String,
    val hostId:String,
    val name:String,
    val gameTemplateId:String,
    val memoryLimitMiB:Int,
    val observedState:ServerState = ServerState.UNKNOWN,
)
