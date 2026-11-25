package io.github.bearbine.easyground

import androidx.compose.runtime.Composable
import io.github.bearbine.easyground.model.GameServer
import androidx.compose.material3.Card
import androidx.compose.material3.Text
import androidx.compose.ui.Modifier
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.padding
import androidx.compose.ui.unit.dp
import easyground.shared.generated.resources.*
import io.github.bearbine.easyground.model.ServerState
import org.jetbrains.compose.resources.stringResource

//=============================
@Suppress("FunctionName")
@Composable

//тестовая карточка сервера при клике (лишь для тестов)
// не пугвйтесь кто попробует запустить
fun ServerCard(
    server:GameServer, //при вызове передаем обьект сервака
    modifier:Modifier = Modifier,//визуальная часть и поведение
    onClick: () -> Unit = {} //передача callback при клике
    // само по себе имеет null    
    ){
      
    // захардкоженный строковой ресурс состояния 
    // не забыть при добавлении нового обновлять
    // или переработать это        
    val stateResource = when(server.observedState){
            ServerState.UNKNOWN -> Res.string.state_unknown
            ServerState.STOPPED -> Res.string.state_stopped
            ServerState.STARTING -> Res.string.state_starting
            ServerState.RUNNING -> Res.string.state_running
            ServerState.STOPPING -> Res.string.state_stopping
    }
    
    //карточка
    Card(
        onClick = onClick,//клик
        modifier = modifier,//визуал
        ){
        Column(modifier = Modifier.padding(16.dp)){
                Text(text = server.name)
                Text(
                    text = stringResource(
                        Res.string.server_memory,
                        server.memoryLimitMiB,
                    ),
                )
                Text(
                    text = stringResource(
                        Res.string.server_state,
                        stringResource(stateResource),
                    ),
                )
                
        }
    }
}