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
@Composable
fun ServerFormScreen(
    onBack:()->Unit,
    onSave:(String, Int, String)->Unit,
){
    //================
    var serverName by rememberSaveable{mutableStateOf("")}
    var memoryInput by rememberSaveable{mutableStateOf("2048")}
    var selectedTemplateId by rememberSaveable{mutableStateOf("game-server")}
    var templateMenuExpanded by remember{mutableStateOf(false)}

    //=========================
    val nameInvalid = serverName.isBlank()
    val memoryMiB = memoryInput.trim().toIntOrNull()
    val memoryInvalid = memoryMiB == null || memoryMiB <= 0
    val formValid = !nameInvalid && !memoryInvalid

    val templates = listOf(
        "game-server" to Res.string.category_game,
        "web-server" to Res.string.category_web,
        "database-server" to Res.string.category_database,
    )
    val selectedTemplateTitle = templates.first{
        it.first == selectedTemplateId
    }.second

    //=========================
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
        Button(
            enabled = formValid,
            onClick = {
                if(formValid && memoryMiB != null){
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

        OutlinedButton(
            onClick = onBack,
            modifier = Modifier.fillMaxWidth(),
        ){
            Text(stringResource(Res.string.cancel))
        }
    }
}