import { Variants } from "framer-motion";

export const bubbleVariant: Variants = {
    hidden: {
        opacity: 0,
        y: 20,
        scale: 0.95,
    },
    visible: {
        opacity: 1,
        y: 0,
        scale: 1,
        transition: {
            type: "spring",
            stiffness: 300,
            damping: 24,
        }
    },
    hover: {
        scale: 1.01,
        transition: {
            type: "spring",
            stiffness: 400,
            damping: 10,
        }
    },
    tap: {
        scale: 0.98,
    }
};

export const dockVariant: Variants = {
    idle: {
        scale: 1,
        y: 0,
    },
    hover: (color: string) => ({
        scale: 1.5,
        y: -10,
        position: "relative",
        zIndex: "111111",
        backgroundColor: color || "#ffffff",
        transition: {
            type: "spring",
            stiffness: 400,
            damping: 25,
        }
    }),
    tap: {
        scale: 1.2,
        y: -5,
    }
};