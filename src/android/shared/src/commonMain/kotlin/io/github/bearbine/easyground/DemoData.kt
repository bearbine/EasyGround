package io.github.bearbine.easyground.demo

import io.github.bearbine.easyground.model.Host
import io.github.bearbine.easyground.model.GameServer
import io.github.bearbine.easyground.model.ServerState

//данные для теста
internal val demoHost = Host(
    id = "demoHost",
    name = "demo host",
)

internal val demoServers = List(20){ index ->
    GameServer(
        id ="demo-server-${index+1}",
        hostId = demoHost.id,
        name = "Demo ${index+1}",
        gameTemplateId = when (index%3){
            0 -> "game-server"
            1 -> "web-server"
            else -> "database-server"
        },
        memoryLimitMiB = 2048,
        observedState = if(index%2==0){
            ServerState.RUNNING
        }else{
            ServerState.STOPPED
        },
    )
}