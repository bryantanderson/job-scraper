const copyFields = (target: object, source: object) => {
    const targetKeys = Object.keys(target);
    targetKeys.forEach((key) => {
        if (key in source) {
            // Typecast to never as we do not need to know the specific type of the object
            (target as never)[key] = (source as never)[key];
        }
    });
};

const wrapText = (text: string) => {
    return {
        message: `${text}`,
    };
};

export { copyFields, wrapText };
