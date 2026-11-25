package io.github.bearbine.easyground

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.runtime.saveable.rememberSaveable
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import easyground.shared.generated.resources.*
import org.jetbrains.compose.resources.stringResource

//=========================
@Suppress("FunctionName")
@Composable

//сейв данных при повороте екрана
fun ServerFormScreen(
    onBack:()->Unit,//callback для выхода из формы
    onSave:(String, Int, String)->Unit,//callback при сохранении
){
    //================
    //хардкоженные вары, потом поменяю    
    var serverName by rememberSaveable{mutableStateOf("")}
    var memoryInput by rememberSaveable{mutableStateOf("2048")}
    var selectedTemplateId by rememberSaveable{mutableStateOf("game-server")}
    var templateMenuExpanded by remember{mutableStateOf(false)}

    //=========================    
    val nameInvalid = serverName.isBlank()//пусто ли
    val memoryMiB = memoryInput.trim().toIntOrNull()//преобразовка в инт
    val memoryInvalid = memoryMiB == null || memoryMiB <= 0//если кто то вставит минусовые значения
    val formValid = !nameInvalid && !memoryInvalid // зеленый свет

    val templates = listOf( //доступные хардкоженные шаблоны
        "game-server" to Res.string.category_game,
        "web-server" to Res.string.category_web,
        "database-server" to Res.string.category_database,
    )    
    //поиск пары по выбранному ид и побране ресурс нейма    
    val selectedTemplateTitle = templates.first{it.first == selectedTemplateId}.second //
    //=========================
    //вертикально, учитываем тут клавиатуру и прокрутку
    Column(
        modifier = Modifier
            .fillMaxSize()
            .safeContentPadding()
            .imePadding()
            .verticalScroll(rememberScrollState())
            .padding(16.dp),
        verticalArrangement = Arrangement.spacedBy(12.dp),
    ){
        Text(
            text = stringResource(Res.string.add_server),
            style = MaterialTheme.typography.headlineSmall,
        )
        //ввод имени сервера при пустом имени показываем ошибку
        OutlinedTextField(
            value = serverName,
            onValueChange = {serverName = it},
            label = {Text(stringResource(Res.string.server_name_label))},
            singleLine = true,
            isError = nameInvalid,
            supportingText = {
                if(nameInvalid){
                    Text(stringResource(Res.string.server_name_error))
                }
            },
            modifier = Modifier.fillMaxWidth(),
        )

        //ввод лимита памяти
        OutlinedTextField(
            value = memoryInput,
            onValueChange = {memoryInput = it},
            label = {Text(stringResource(Res.string.server_memory_label))},
            singleLine = true,
            isError = memoryInvalid,
            supportingText = {
                if(memoryInvalid){
                    Text(stringResource(Res.string.server_memory_error))
                }
            },
            modifier = Modifier.fillMaxWidth(),
        )

        //=========================
        //кнопка выбранного шаблона        
        Box{
            Button(onClick = {templateMenuExpanded = true}){
                Text(stringResource(selectedTemplateTitle))
            }

            DropdownMenu(
                expanded = templateMenuExpanded,
                onDismissRequest = {templateMenuExpanded = false},
            ){
                templates.forEach{(templateId, title) ->
                    DropdownMenuItem(
                        text = {Text(stringResource(title))},
                        onClick = {
                            selectedTemplateId = templateId
                            templateMenuExpanded = false
                        },
                    )
                }
            }
        }

        //=========================
        //сохранить сервер        
        Button(
            enabled = formValid,
            onClick = {
                if(formValid){
                    onSave(
                        serverName.trim(),
                        memoryMiB,
                        selectedTemplateId,
                    )
                }
            },
            modifier = Modifier.fillMaxWidth(),
        ){
            Text(stringResource(Res.string.save_server))
        }
        //отменить ввод
        OutlinedButton(
            onClick = onBack,
            modifier = Modifier.fillMaxWidth(),
        ){
            Text(stringResource(Res.string.cancel))
        }
    }
}