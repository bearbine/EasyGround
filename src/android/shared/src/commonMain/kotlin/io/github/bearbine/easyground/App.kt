package io.github.bearbine.easyground

import androidx.compose.material3.MaterialTheme
import androidx.compose.runtime.Composable
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.runtime.getValue
import androidx.compose.runtime.setValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.saveable.rememberSaveable
import io.github.bearbine.easyground.demo.demoServers
import androidx.compose.runtime.saveable.rememberSaveableStateHolder

//=========================
@Suppress("FunctionName")
@Composable
@Preview

//выбор какой екран будет на показе у юзера
fun App(){
    MaterialTheme{
        //ид выбранного сервера, список будет нулевый    
        var selectedServerId by rememberSaveable{ mutableStateOf<String?>(null)}
        //===================
        //ищем данные выбранного сервера
        val selectedServer = demoServers.find{it.id == selectedServerId}
        //сохраняем поиск фильтры и прокрутку если перейдем кудато
        val screenStateHolder = rememberSaveableStateHolder()
        //================

        
        if(selectedServer == null){
            screenStateHolder.SaveableStateProvider("server-list"){
                ServerListScreen(
                    onServerClick = {server ->
                    selectedServerId = server.id
                    },
                )
            }
        }else{
            ServerDetailsScreen(
                server = selectedServer,
                onBack = {selectedServerId = null},
            )
        }
    }
}