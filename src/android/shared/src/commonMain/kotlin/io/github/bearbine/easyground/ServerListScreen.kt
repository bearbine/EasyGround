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
import androidx.compose.material3.OutlinedTextField
import androidx.compose.runtime.saveable.rememberSaveable
import androidx.compose.material3.FilterChip
import io.github.bearbine.easyground.model.ServerState
import androidx.compose.foundation.lazy.LazyRow
import io.github.bearbine.easyground.model.GameServer
import easyground.shared.generated.resources.*
    
//=====================
@Suppress("FunctionName")
@Composable

//показ списка с поиском и фильтром
fun ServerListScreen(
    onServerClick: (GameServer) -> Unit = {}, //callback от клика в КОНКРЕТНЫЙ обработчик
){
    
        //==============      
        var searchQuery by rememberSaveable{mutableStateOf("")} //текст в строке
        var runningOnly by rememberSaveable{mutableStateOf(false)} //включен ли фильтр включенных серверов
        var selectedCategory by rememberSaveable{mutableStateOf<String?>(null)} //выбранная ката
        var categories = listOf( //варианты кат
        // да тут хардкод(пока что)        
            null to Res.string.category_all,
            "game-server" to Res.string.category_game,
            "web-server" to Res.string.category_web,
            "database-server" to Res.string.category_database,
        )
        //=====================
        //серверы которые прошли фильтр        
        val filteredServers = demoServers.filter{server ->
            //матч неймов подходит ли         
            val matchesName = server.name.contains(
                searchQuery.trim(),
                ignoreCase = true,
            )
            //подходит ли состояние и категория
            //доп проверки если буду расширять то не забыть сюда закинуть                        
            val matchesState = !runningOnly || server.observedState == ServerState.RUNNING
            val matchesCategory = selectedCategory == null || server.gameTemplateId == selectedCategory

            matchesName && matchesState && matchesCategory
        }
        //=====================
        //вертикальный список        
        LazyColumn(
        modifier = Modifier.fillMaxSize().safeContentPadding(),//занимаем все место учитывая отступы
        contentPadding = PaddingValues(16.dp), //отступ содержимого
        verticalArrangement = Arrangement.spacedBy(12.dp),//лонг между елементами
       ){

        //демо данные
        item{
            Text(text = stringResource(Res.string.demo_mode))
        }
                
        //поле поиска
        item{
            OutlinedTextField(
                value = searchQuery,
                onValueChange = {searchQuery = it},
                label = {
                    Text(stringResource(Res.string.search_servers))
                },
                singleLine = true,
                modifier = Modifier.fillMaxWidth()
            )
        }

        //выбор категории
        item{
            LazyRow(
                horizontalArrangement = Arrangement.spacedBy(8.dp),
            ){
                items(categories){(categoryId, title) ->
                    FilterChip(
                        selected = selectedCategory == categoryId,
                        onClick = {selectedCategory = categoryId},
                        label = {Text(stringResource(title))},
                    )
                }
            }
        }

        //фильтр только запущенных
        item{
            FilterChip(
                selected = runningOnly,
                onClick = { runningOnly = !runningOnly},
                label = {
                    Text(stringResource(Res.string.running_only))
                },
            )
        }
        
        //сообщение если список будет пустым                
        if(filteredServers.isEmpty()){
            item{
                Text(
                    text = stringResource(
                        Res.string.no_matching_servers,
                    ),
                )
            }
        }

        //карточки тех кто выжил и прошел фильтры
        items(filteredServers, key = {it.id}){ server ->
            ServerCard(
                server = server,
                modifier = Modifier.fillMaxWidth(),
                onClick = {onServerClick(server)},
            )
        }
    }
}