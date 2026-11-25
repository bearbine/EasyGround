package io.github.bearbine.easyground

import androidx.compose.animation.AnimatedVisibility
import androidx.compose.foundation.Image
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.PaddingValues
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.ui.unit.dp
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.safeContentPadding
import androidx.compose.material3.Button
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.tooling.preview.Preview
import io.github.bearbine.easyground.demo.demoServers
import easyground.shared.generated.resources.demo_mode
import org.jetbrains.compose.resources.stringResource
import org.jetbrains.compose.resources.painterResource

import easyground.shared.generated.resources.*

@Composable
@Preview
fun App(){
    MaterialTheme {
       LazyColumn(
        modifier = Modifier.fillMaxSize().safeContentPadding(),
        contentPadding = PaddingValues(16.dp),
        verticalArrangement = Arrangement.spacedBy(12.dp),
       ){
        item{
            Text(text = stringResource(Res.string.demo_mode))
        }

        items(demoServers, key = {it.id}){ server ->
            ServerCard(
                server = server,
                modifier = Modifier.fillMaxWidth(),
            )
        }
       }
    }
}