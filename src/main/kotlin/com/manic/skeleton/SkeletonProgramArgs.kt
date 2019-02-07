package com.manic.skeleton

import com.xenomachina.argparser.ArgParser

/**
 * Data class for the __skeleton__ program.
 */
class SkeletonProgramArgs(parser: ArgParser) {
    val command: String by parser.positional("CMD")
}