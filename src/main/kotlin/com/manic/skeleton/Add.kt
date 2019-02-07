package com.manic.skeleton

import com.xenomachina.argparser.SystemExitException
import java.io.File
import java.io.FileWriter

fun add(addCommandArgs: AddCommandArgs) {
    val domainName = addCommandArgs.domainName
    val parent = File(domainName)

    if (parent.exists()) {
        throw SystemExitException("Directory ${parent.absoluteFile} is already existing.", 101)
    }

    val application = File("$domainName/$domainName-application")
    val domain = File("$domainName/$domainName-domain")
    val infrastructure = File("$domainName/$domainName-infrastructure")

    arrayOf(
        parent,
        application,
        domain,
        infrastructure
    ).forEach { it.mkdir() }

    val groupId = read("group id")
    val artifactId = read("artifact id", domainName)
    val version = read("version")

    val parentPom = FreemarkerConfig.getTemplate("add/parent-pom.ftlh")
    parentPom.process(
        mapOf(
            "groupId" to groupId,
            "artifactId" to artifactId,
            "version" to version
        ), FileWriter(File(parent, "pom.xml"))
    )

    println("Domain created!")
}

private fun read(propertyName: String, suggestion: String? = null): String? {
    if (suggestion != null) {
        print("Take $propertyName '$suggestion'?[y/n]: ")
        when (readLine()) {
            "y" -> return suggestion
        }
    }

    var property: String?
    do {
        print("Type $propertyName: ")
        property = readLine() ?: ""
        print("${propertyName.capitalize()} '$property' correct?[y/n]: ")
    } while (readLine() != "y")
    return property
}