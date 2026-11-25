package io.github.bearbine.easyground.model

//последнее известное состояние игрового сервера
enum class ServerState{
    UNKNOWN,
    STOPPED,
    STARTING,
    RUNNING,
    STOPPING,
}