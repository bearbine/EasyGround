package io.github.bearbine.easyground.demo

import io.github.bearbine.easyground.model.Host
import io.github.bearbine.easyground.model.GameServer
import io.github.bearbine.easyground.model.ServerState

//====================
//данные для теста
//демохост по аналогичному принципу будет создаваться основные
internal val demoHost = Host(
    id = "demoHost", 
    name = "demo host",
)

//создаю демосервер
//index идет включая 0 до 19, создавай 20 демохостов
//хардкодом легче это описать
internal val demoServers = List(20){ index ->
    GameServer(
        ////string template в формате ${index+X}    
        id ="demo-server-${index+1}",
        hostId = demoHost.id,
        name = "Demo ${index+1}",
        
        //для создания масовки даем псевдорандомно по модуль3 темплейт                
        gameTemplateId = when (index%3){
            0 -> "game-server"
            1 -> "web-server"
            else -> "database-server"
        },
        memoryLimitMiB = 2048,
        //для создания масовки даем по псевдорандому состояния        
        observedState = if(index%2==0){
            ServerState.RUNNING
        }else{
            ServerState.STOPPED
        },
    )
}