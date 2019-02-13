package com.manic.skeleton

import com.xenomachina.argparser.SystemExitException
import kotlinx.serialization.Serializable
import kotlinx.serialization.json.Json.Companion.parse
import java.io.File
import java.io.FileWriter

@Serializable
data class ProjectDescriptor(
    val version: String,
    val name: String,
    val description: String?,
    val variables: Array<String>,
    val structure: Array<Structure>
)

@Serializable
data class Structure(
    val dir: String,
    val file: String,
    val template: String
)

fun add(addCommandArgs: AddCommandArgs) {
    val projectDescriptor = parse(
        ProjectDescriptor.serializer(),
        ClassLoader.getSystemClassLoader().getResourceAsStream("templates/add/maven-domain-module.json").reader().readText()
    )

    val readVariablesMap = read(projectDescriptor.variables)
    projectDescriptor.structure.forEach { structure ->
        val actualDir = replace(structure.dir, readVariablesMap)
        val actualFile = File(actualDir)
        actualFile.mkdir()

        FreemarkerConfig.getTemplate("add/${structure.template}").process(
            readVariablesMap,
            FileWriter(File(actualFile, structure.file))
        )
    }
    println("Domain created!")
}

private fun replace(
    value: String,
    readVariablesMap: Map<String, String?>
): String {
    return value.replace(Regex("\\$\\{(.*?)}")) {
        readVariablesMap[it.destructured.component1()] ?: throw SystemExitException(
            "Couldnt find value for variable '${it.value}'",
            101
        )
    }
}

private fun read(propertyNames: Array<String>): Map<String, String?> {
    val readProperties = mutableMapOf<String, String?>()
    do {
        readProperties.clear()
        propertyNames.forEach { propertyName ->
            print("Type $propertyName: ")
            readProperties[propertyName] = readLine() ?: ""
        }

        println()
        println("{")
        readProperties.forEach { println("\"${it.key}\": \"${it.value}\"") }
        println("}")
        println()

        print("Are values correct?[y/n]: ")
    } while (readLine() != "y")
    return readProperties
}