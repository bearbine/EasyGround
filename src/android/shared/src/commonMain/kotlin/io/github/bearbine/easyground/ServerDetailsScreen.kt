package io.github.bearbine.easyground

import androidx.compose.foundation.layout.Column
import androidx.compose.material3.Button
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import io.github.bearbine.easyground.model.GameServer
import easyground.shared.generated.resources.*
import org.jetbrains.compose.resources.stringResource
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.safeContentPadding
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import io.github.bearbine.easyground.model.ServerState

//====================
@Composable
fun ServerDetailsScreen(
    server:GameServer,
    onBack:()->Unit,
){
    //=============
    val stateResource = when (server.observedState) {
        ServerState.UNKNOWN -> Res.string.state_unknown
        ServerState.STOPPED -> Res.string.state_stopped
        ServerState.STARTING -> Res.string.state_starting
        ServerState.RUNNING -> Res.string.state_running
        ServerState.STOPPING -> Res.string.state_stopping
    }
    //==========================
    Column(
        modifier = Modifier
        .fillMaxSize()
        .safeContentPadding()
        .verticalScroll(rememberScrollState())
        .padding(16.dp),
        verticalArrangement = Arrangement.spacedBy(12.dp),
    ){
        Text(server.name)
        Text(stringResource(Res.string.server_id, server.id))
        Text(stringResource(Res.string.server_host_id, server.hostId))
        Text(stringResource(Res.string.server_template, server.gameTemplateId))
        Text(stringResource(Res.string.server_memory, server.memoryLimitMiB))
        Text(stringResource(Res.string.server_state,stringResource(stateResource),),)
        Button(onClick = onBack){
            Text(stringResource(Res.string.back_to_servers))
        }
    }
}