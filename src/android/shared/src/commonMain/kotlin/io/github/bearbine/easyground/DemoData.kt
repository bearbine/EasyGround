package io.github.bearbine.easyground.demo

import io.github.bearbine.easyground.model.Host
import io.github.bearbine.easyground.model.GameServer
import io.github.bearbine.easyground.model.ServerState

//данные для теста
internal val demoHost = Host(
    id = "demoHost",
    name = "demo host",
)

internal val demoServers = listOf(
    GameServer(
        id ="demo-server",
        hostId = demoHost.id,
        name = "Demo-name",
        gameTemplateId = "server jar",
        memoryLimitMiB = 2048,
        observedState = ServerState.RUNNING,
    ),
)