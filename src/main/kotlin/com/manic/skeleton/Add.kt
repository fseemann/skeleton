package com.manic.skeleton

import com.xenomachina.argparser.SystemExitException
import java.io.File
import java.io.FileWriter

fun add(addCommandArgs: AddCommandArgs) {
    val domainName = addCommandArgs.domainName

    val groupId = read("group id")
    val artifactId = read("artifact id", domainName)
    val version = read("version")

    val parent = File(artifactId)
    if (parent.exists()) {
        throw SystemExitException("Directory ${parent.absoluteFile} is already existing.", 101)
    }

    val application = File("$artifactId/$artifactId-application")
    val domain = File("$artifactId/$artifactId-domain")
    val infrastructure = File("$artifactId/$artifactId-infrastructure")

    arrayOf(
        parent,
        application,
        domain,
        infrastructure
    ).forEach { it.mkdir() }

    listOf(
        Triple(parent, "add/parent-pom.ftlh", "pom.xml"),
        Triple(application, "add/application-pom.ftlh", "pom.xml"),
        Triple(domain, "add/domain-pom.ftlh", "pom.xml"),
        Triple(infrastructure, "add/infrastructure-pom.ftlh", "pom.xml")
    ).forEach {
        FreemarkerConfig.getTemplate(it.second).process(
            mapOf(
                "groupId" to groupId,
                "artifactId" to artifactId,
                "version" to version
            ), FileWriter(File(it.first, it.third))
        )
    }
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