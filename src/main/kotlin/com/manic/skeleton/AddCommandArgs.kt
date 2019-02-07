package com.manic.skeleton

import com.xenomachina.argparser.ArgParser

/**
 * Data class for the skeleton __add__ command.
 */
class AddCommandArgs(parser: ArgParser) {
    val domainName: String by parser.positional("DOMAIN_NAME")
}